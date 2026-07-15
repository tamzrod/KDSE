#!/usr/bin/env python3
"""
KDSE Slim Laboratory Runner
Executes laboratory scenarios to validate KDSE methodology.
"""

import sys
import os
import json
import argparse

# Add laboratory core to path
sys.path.insert(0, os.path.join(os.path.dirname(__file__), 'core'))

from core.laboratory import (
    LaboratoryRunner, 
    LaboratoryScenario,
    create_lab001_scenario,
    LaboratoryResult
)


def run_scenario(scenario_id: str, workspace_root: str = ".kdse", 
                 output_dir: str = "results", objective: str = None) -> bool:
    """
    Run a laboratory scenario.
    
    Args:
        scenario_id: ID of scenario to run (e.g., "LAB-001")
        workspace_root: Root directory for KDSE workspace
        output_dir: Directory to write results
        objective: Custom objective for the scenario
        
    Returns:
        True if experiment passed, False otherwise
    """
    print(f"KDSE Slim Laboratory Runner")
    print(f"Scenario: {scenario_id}")
    print(f"Workspace: {workspace_root}")
    if objective:
        print(f"Objective: {objective}")
    print("=" * 60)
    
    # Create scenario based on ID
    if scenario_id == "LAB-001":
        scenario = create_lab001_scenario(objective)
    else:
        print(f"Unknown scenario: {scenario_id}")
        return False
    
    # Create workspace directory
    os.makedirs(workspace_root, exist_ok=True)
    os.makedirs(os.path.join(workspace_root, "knowledge"), exist_ok=True)
    os.makedirs(output_dir, exist_ok=True)
    
    # Run the laboratory
    runner = LaboratoryRunner(scenario, workspace_root)
    report = runner.run()
    
    # Save report
    report_path = os.path.join(output_dir, f"{scenario_id}-REPORT.json")
    with open(report_path, 'w') as f:
        json.dump(report.to_dict(), f, indent=2)
    
    # Generate markdown report
    md_report = generate_markdown_report(report)
    md_path = os.path.join(output_dir, f"{scenario_id}-REPORT.md")
    with open(md_path, 'w') as f:
        f.write(md_report)
    
    print(f"\nReports saved to:")
    print(f"  JSON: {report_path}")
    print(f"  Markdown: {md_path}")
    
    return report.result == LaboratoryResult.PASS


def generate_markdown_report(report) -> str:
    """Generate a markdown report from the laboratory report."""
    
    status_icon = "✅ PASS" if report.result == LaboratoryResult.PASS else "❌ FAIL"
    
    md = f"""# {report.scenario_id} Laboratory Report

## Experiment Result: {status_icon}

**Scenario:** {report.scenario_title}  
**Objective:** {report.objective}  
**Started:** {report.started_at}  
**Completed:** {report.completed_at}

---

## Summary

{report.pass_reason if report.result == LaboratoryResult.PASS else report.fail_reason}

---

## Steps Executed

"""
    
    for step in report.steps_executed:
        md += f"""### {step['step'].replace('_', ' ').title()}

"""
        if 'artifacts' in step:
            for artifact in step['artifacts']:
                md += f"- Created: `{artifact.get('name', 'unknown')}`\n"
        
        if 'facts_identified' in step:
            md += f"- Facts identified: {step.get('facts_identified', 0)}\n"
        
        if 'assumptions_identified' in step:
            md += f"- Assumptions identified: {step.get('assumptions_identified', 0)}\n"
        
        if 'gaps_identified' in step:
            md += f"- Gaps identified: {step.get('gaps_identified', 0)}\n"
        
        if 'summary' in step:
            md += f"- Evidence Strength: {step['summary']}\n"
        
        md += "\n"
    
    md += """---

## Knowledge Distinction

### Known Facts
"""
    
    for fact in report.knowledge_distinction.get('known_facts', []):
        md += f"""**{fact.get('statement', 'Unknown')}**
- Evidence: {fact.get('evidence', 'N/A')}
- Strength: {fact.get('strength', '●')}

"""
    
    md += """### Reasonable Assumptions
"""
    
    for assumption in report.knowledge_distinction.get('assumptions', []):
        md += f"""**{assumption.get('statement', 'Unknown')}**
- Basis: {assumption.get('basis', 'N/A')}
- Status: {assumption.get('status', 'assumed')}
- Strength: {assumption.get('strength', '●')}

"""
    
    md += """### Knowledge Gaps
"""
    
    for gap in report.knowledge_distinction.get('knowledge_gaps', []):
        md += f"""**{gap.get('category', 'unknown').title()}: {gap.get('statement', 'Unknown')}**
- Investigation Needed: {gap.get('investigation', 'N/A')}

"""
    
    md += """---

## Evidence Strength Summary

"""
    
    strength = report.knowledge_distinction.get('evidence_strength', {})
    md += f"""- Total Statements: {strength.get('total_statements', 0)}
- ●●● (Strong): {strength.get('strong', 0)}
- ●● (Moderate): {strength.get('moderate', 0)}
- ● (Weak): {strength.get('weak', 0)}

"""
    
    md += """---

## Foundation Readiness

"""
    
    foundation_checkpoints = [
        ("foundation_skeleton_created", "Foundation Skeleton Created"),
        ("facts_identified", "Known Facts Identified"),
        ("assumptions_distinguished", "Assumptions Distinguished"),
        ("gaps_identified", "Knowledge Gaps Identified"),
        ("evidence_strength_recorded", "Evidence Strength Recorded"),
        ("foundation_refinement_planned", "Refinement Planned")
    ]
    
    for checkpoint_id, description in foundation_checkpoints:
        status = "✅" if checkpoint_id in report.checkpoints_passed else "❌"
        md += f"- {status} {description}\n"
    
    md += """
---

## Next Decisions

"""
    
    if report.decisions_made:
        for i, decision in enumerate(report.decisions_made[:3], 1):
            md += f"""### {i}. {decision.get('objective', 'Unknown')}

- **Type:** {decision.get('decision_type', 'N/A')}
- **Priority:** {decision.get('priority', 'N/A')}
- **Traces To:** {', '.join(decision.get('traces_to', ['None'])) or 'None'}

"""
    else:
        md += "No additional decisions required.\n"
    
    md += """---

## Findings

"""
    
    for finding in report.findings:
        md += f"- {finding}\n"
    
    md += """
---

## Artifacts Created

"""
    
    for artifact in report.artifacts_created:
        md += f"- `{artifact}`\n"
    
    if report.forbidden_actions:
        md += """
---

## Forbidden Behavior Detected ⚠️

"""
        for action in report.forbidden_actions:
            md += f"- ❌ {action}\n"
    
    md += f"""

---

## Final Verdict

**{"✅ PASS" if report.result == LaboratoryResult.PASS else "❌ FAIL"}**: {report.pass_reason if report.result == LaboratoryResult.PASS else report.fail_reason}

---

*Report generated: {report.completed_at}*
"""
    
    return md


def main():
    parser = argparse.ArgumentParser(description='KDSE Slim Laboratory Runner')
    parser.add_argument('--scenario', default='LAB-001', help='Scenario ID to run')
    parser.add_argument('--workspace', default='.kdse', help='KDSE workspace root')
    parser.add_argument('--output', default='results', help='Output directory for reports')
    parser.add_argument('--objective', default=None, help='Custom objective for the scenario')
    
    args = parser.parse_args()
    
    success = run_scenario(args.scenario, args.workspace, args.output, args.objective)
    sys.exit(0 if success else 1)


if __name__ == '__main__':
    main()

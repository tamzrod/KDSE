"""
KDSE Slim - Derivation Orchestrator
Guides the GATHER → DERIVE → VALIDATE process.
"""

from dataclasses import dataclass, field
from enum import Enum
from typing import List, Optional, Dict
from datetime import datetime
import uuid


class DerivationStage(Enum):
    """Stages of the derivation process."""
    GATHER = "gather"      # Find and classify evidence
    DERIVE = "derive"      # Transform evidence into knowledge
    VALIDATE = "validate"  # Check derivation quality


class EvidenceStrength(Enum):
    """Simplified evidence strength scale."""
    WEAK = "●"          # Single source or inference
    MODERATE = "●●"     # Some corroboration
    STRONG = "●●●"      # Multiple independent sources


@dataclass
class EvidenceSource:
    """A source of evidence."""
    id: str
    type: str  # code, document, standard, vendor, etc.
    location: str
    finding: str  # What this evidence shows
    provenance: str = ""  # Origin information
    timestamp: str = ""
    
    def __post_init__(self):
        if not self.timestamp:
            self.timestamp = datetime.now().isoformat()


@dataclass
class KnowledgeStatement:
    """A single knowledge statement."""
    id: str
    statement: str  # The knowledge claim
    evidence_ids: List[str] = field(default_factory=list)
    strength: EvidenceStrength = EvidenceStrength.WEAK
    category: str = ""  # requirement, constraint, assumption, context, etc.
    traces_to: List[str] = field(default_factory=list)  # Other knowledge this depends on
    derived_by: str = ""
    derived_at: str = ""
    validated: bool = False
    
    def __post_init__(self):
        if not self.id:
            self.id = f"KNOW-{uuid.uuid4().hex[:8].upper()}"
        if not self.derived_at:
            self.derived_at = datetime.now().isoformat()


class DerivationOrchestrator:
    """
    Guides the GATHER → DERIVE → VALIDATE process.
    
    This is NOT a phase runner. It guides reasoning.
    """
    
    def __init__(self):
        self.evidence_sources: List[EvidenceSource] = []
        self.knowledge_statements: List[KnowledgeStatement] = []
        self.current_stage: DerivationStage = DerivationStage.GATHER
    
    def reset(self):
        """Reset to initial state."""
        self.evidence_sources = []
        self.knowledge_statements = []
        self.current_stage = DerivationStage.GATHER
    
    # ========== GATHER Stage ==========
    
    def gather(self, evidence: List[Dict]) -> List[EvidenceSource]:
        """
        GATHER Stage: Find and classify evidence sources.
        
        Args:
            evidence: List of evidence dicts with type, location, finding
            
        Returns:
            List of EvidenceSource objects
        """
        self.current_stage = DerivationStage.GATHER
        
        sources = []
        for e in evidence:
            source = EvidenceSource(
                id=f"EVID-{uuid.uuid4().hex[:8].upper()}",
                type=e.get("type", "unknown"),
                location=e.get("location", ""),
                finding=e.get("finding", ""),
                provenance=e.get("provenance", "")
            )
            sources.append(source)
            self.evidence_sources.append(source)
        
        return sources
    
    def add_evidence(self, evidence_type: str, location: str, finding: str) -> EvidenceSource:
        """Add a single evidence source."""
        source = EvidenceSource(
            id=f"EVID-{uuid.uuid4().hex[:8].upper()}",
            type=evidence_type,
            location=location,
            finding=finding
        )
        self.evidence_sources.append(source)
        return source
    
    def get_evidence_for_finding(self, finding: str) -> List[EvidenceSource]:
        """Get evidence sources related to a finding."""
        finding_lower = finding.lower()
        return [
            e for e in self.evidence_sources
            if finding_lower in e.finding.lower() or e.finding.lower() in finding_lower
        ]
    
    # ========== DERIVE Stage ==========
    
    def derive(self, statements: List[Dict]) -> List[KnowledgeStatement]:
        """
        DERIVE Stage: Transform evidence into knowledge statements.
        
        This is NOT extraction. It requires reasoning.
        
        Args:
            statements: List of knowledge claim dicts
            
        Returns:
            List of KnowledgeStatement objects
        """
        self.current_stage = DerivationStage.DERIVE
        
        knowledge_items = []
        for s in statements:
            # Find supporting evidence
            evidence_ids = []
            for e in self.evidence_sources:
                if s.get("evidence_supports", "").lower() in e.finding.lower():
                    evidence_ids.append(e.id)
            
            # Calculate strength based on evidence count
            strength = EvidenceStrength.WEAK
            if len(evidence_ids) >= 3:
                strength = EvidenceStrength.STRONG
            elif len(evidence_ids) >= 2:
                strength = EvidenceStrength.MODERATE
            
            statement = KnowledgeStatement(
                id=s.get("id", f"KNOW-{uuid.uuid4().hex[:8].upper()}"),
                statement=s.get("statement", ""),
                evidence_ids=evidence_ids,
                strength=strength,
                category=s.get("category", "general"),
                traces_to=s.get("traces_to", []),
                derived_by=s.get("derived_by", ""),
                validated=False
            )
            knowledge_items.append(statement)
            self.knowledge_statements.append(statement)
        
        return knowledge_items
    
    def derive_from_evidence(self, evidence_id: str, claim: str, category: str = "general") -> KnowledgeStatement:
        """
        Derive a single knowledge statement from evidence.
        """
        statement = KnowledgeStatement(
            id=f"KNOW-{uuid.uuid4().hex[:8].upper()}",
            statement=claim,
            evidence_ids=[evidence_id],
            strength=EvidenceStrength.MODERATE if evidence_id else EvidenceStrength.WEAK,
            category=category,
            derived_by="derivation_orchestrator",
            validated=False
        )
        self.knowledge_statements.append(statement)
        return statement
    
    def get_underived_evidence(self) -> List[EvidenceSource]:
        """Get evidence that hasn't been used to derive knowledge."""
        used_evidence = set()
        for k in self.knowledge_statements:
            used_evidence.update(k.evidence_ids)
        return [e for e in self.evidence_sources if e.id not in used_evidence]
    
    # ========== VALIDATE Stage ==========
    
    def validate(self) -> Dict:
        """
        VALIDATE Stage: Check derivation quality and assign final strength.
        
        Returns validation report.
        """
        self.current_stage = DerivationStage.VALIDATE
        
        validation_results = {
            "total_statements": len(self.knowledge_statements),
            "statements_with_evidence": 0,
            "statements_without_evidence": 0,
            "weak_statements": 0,
            "moderate_statements": 0,
            "strong_statements": 0,
            "unvalidated_statements": [],
            "validation_timestamp": datetime.now().isoformat()
        }
        
        for statement in self.knowledge_statements:
            # Check if has evidence
            if statement.evidence_ids:
                validation_results["statements_with_evidence"] += 1
            else:
                validation_results["statements_without_evidence"] += 1
                validation_results["unvalidated_statements"].append(statement.id)
            
            # Count by strength
            if statement.strength == EvidenceStrength.WEAK:
                validation_results["weak_statements"] += 1
            elif statement.strength == EvidenceStrength.MODERATE:
                validation_results["moderate_statements"] += 1
            elif statement.strength == EvidenceStrength.STRONG:
                validation_results["strong_statements"] += 1
        
        # Mark validated
        for statement in self.knowledge_statements:
            statement.validated = True
        
        return validation_results
    
    def assign_strength(self, statement_id: str, strength: EvidenceStrength):
        """Manually assign evidence strength to a statement."""
        for statement in self.knowledge_statements:
            if statement.id == statement_id:
                statement.strength = strength
                break
    
    # ========== Full Cycle ==========
    
    def run_derivation_cycle(self, evidence: List[Dict], claims: List[Dict]) -> Dict:
        """
        Run the full GATHER → DERIVE → VALIDATE cycle.
        
        Returns comprehensive derivation report.
        """
        # Stage 1: GATHER
        gathered_evidence = self.gather(evidence)
        
        # Stage 2: DERIVE
        derived_knowledge = self.derive(claims)
        
        # Stage 3: VALIDATE
        validation = self.validate()
        
        return {
            "stage": self.current_stage.value,
            "evidence_gathered": len(gathered_evidence),
            "knowledge_derived": len(derived_knowledge),
            "validation": validation,
            "knowledge": [
                {
                    "id": k.id,
                    "statement": k.statement,
                    "strength": k.strength.value,
                    "category": k.category,
                    "evidence_count": len(k.evidence_ids),
                    "validated": k.validated
                }
                for k in self.knowledge_statements
            ],
            "underived_evidence": [
                {"id": e.id, "finding": e.finding}
                for e in self.get_underived_evidence()
            ]
        }
    
    def export_state(self) -> Dict:
        """Export derivation state for persistence."""
        return {
            "evidence_sources": [
                {
                    "id": e.id,
                    "type": e.type,
                    "location": e.location,
                    "finding": e.finding,
                    "provenance": e.provenance,
                    "timestamp": e.timestamp
                }
                for e in self.evidence_sources
            ],
            "knowledge_statements": [
                {
                    "id": k.id,
                    "statement": k.statement,
                    "evidence_ids": k.evidence_ids,
                    "strength": k.strength.value,
                    "category": k.category,
                    "traces_to": k.traces_to,
                    "derived_by": k.derived_by,
                    "derived_at": k.derived_at,
                    "validated": k.validated
                }
                for k in self.knowledge_statements
            ],
            "current_stage": self.current_stage.value
        }

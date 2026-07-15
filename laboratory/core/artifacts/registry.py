"""
KDSE Slim - Artifact Registry
Tracks all engineering artifacts by type, state, and owner.
"""

from dataclasses import dataclass, field
from enum import Enum
from typing import Optional, List
from datetime import datetime


class ArtifactType(Enum):
    """Artifact types organized by authority level."""
    KNOWLEDGE = "knowledge"           # Highest authority
    ARCHITECTURE = "architecture"     # Traces to knowledge
    IMPLEMENTATION = "implementation" # Traces to architecture
    VERIFICATION = "verification"    # Confirms alignment
    EVIDENCE = "evidence"            # Evidence sources


class LifecycleState(Enum):
    """Simplified artifact lifecycle states."""
    DRAFT = "draft"        # Initial creation
    REVIEWED = "reviewed"  # Peer reviewed
    APPROVED = "approved"  # Authoritative
    ARCHIVED = "archived"  # Superseded


@dataclass
class Artifact:
    """A single engineering artifact."""
    id: str
    name: str
    type: ArtifactType
    state: LifecycleState = LifecycleState.DRAFT
    owner: str = ""
    created_at: str = ""
    updated_at: str = ""
    traces_to: List[str] = field(default_factory=list)  # Artifact IDs this traces to
    evidence_strength: str = ""  # ● (weak), ●● (moderate), ●●● (strong)
    description: str = ""
    location: str = ""  # File path within workspace
    
    def __post_init__(self):
        if not self.created_at:
            self.created_at = datetime.now().isoformat()
        if not self.updated_at:
            self.updated_at = self.created_at
    
    def update_state(self, new_state: LifecycleState):
        """Update artifact lifecycle state."""
        self.state = new_state
        self.updated_at = datetime.now().isoformat()
    
    def add_trace(self, target_id: str):
        """Add a trace to another artifact."""
        if target_id not in self.traces_to:
            self.traces_to.append(target_id)
            self.updated_at = datetime.now().isoformat()
    
    def get_authority_level(self) -> int:
        """Get numeric authority level (higher = more authoritative)."""
        authority_map = {
            LifecycleState.DRAFT: 1,
            LifecycleState.REVIEWED: 2,
            LifecycleState.APPROVED: 3,
            LifecycleState.ARCHIVED: 0
        }
        return authority_map.get(self.state, 0)


class ArtifactRegistry:
    """
    Central registry for all KDSE artifacts.
    Tracks what exists, its type, state, owner, and traces.
    """
    
    def __init__(self):
        self.artifacts: dict[str, Artifact] = {}
    
    def register(self, artifact: Artifact) -> None:
        """Register a new artifact."""
        self.artifacts[artifact.id] = artifact
    
    def get(self, artifact_id: str) -> Optional[Artifact]:
        """Get artifact by ID."""
        return self.artifacts.get(artifact_id)
    
    def list_by_type(self, artifact_type: ArtifactType) -> List[Artifact]:
        """List all artifacts of a given type."""
        return [a for a in self.artifacts.values() if a.type == artifact_type]
    
    def list_by_state(self, state: LifecycleState) -> List[Artifact]:
        """List all artifacts in a given state."""
        return [a for a in self.artifacts.values() if a.state == state]
    
    def list_approved(self) -> List[Artifact]:
        """List all approved artifacts (authoritative)."""
        return self.list_by_state(LifecycleState.APPROVED)
    
    def list_drafts(self) -> List[Artifact]:
        """List all draft artifacts."""
        return self.list_by_state(LifecycleState.DRAFT)
    
    def get_traces_from(self, artifact_id: str) -> List[Artifact]:
        """Get all artifacts that the given artifact traces to."""
        artifact = self.get(artifact_id)
        if not artifact:
            return []
        return [self.get(tid) for tid in artifact.traces_to if self.get(tid)]
    
    def get_traces_to(self, artifact_id: str) -> List[Artifact]:
        """Get all artifacts that trace to the given artifact."""
        return [a for a in self.artifacts.values() if artifact_id in a.traces_to]
    
    def check_authority_violation(self, artifact_id: str) -> List[str]:
        """
        Check for authority violations.
        Returns list of violation messages.
        """
        violations = []
        artifact = self.get(artifact_id)
        if not artifact:
            return ["Artifact not found"]
        
        for trace_id in artifact.traces_to:
            target = self.get(trace_id)
            if not target:
                violations.append(f"Trace to non-existent artifact: {trace_id}")
                continue
            
            # Lower type cannot trace to higher authority artifact
            type_hierarchy = {
                ArtifactType.IMPLEMENTATION: 1,
                ArtifactType.ARCHITECTURE: 2,
                ArtifactType.KNOWLEDGE: 3,
                ArtifactType.VERIFICATION: 0,
                ArtifactType.EVIDENCE: 0
            }
            
            # If artifact type is lower authority than target, check states
            if type_hierarchy.get(artifact.type, 0) < type_hierarchy.get(target.type, 0):
                if target.state == LifecycleState.DRAFT:
                    violations.append(
                        f"{artifact.type.value} traces to draft {target.type.value}: {target.id}"
                    )
        
        return violations
    
    def get_summary(self) -> dict:
        """Get registry summary."""
        return {
            "total_artifacts": len(self.artifacts),
            "by_type": {
                t.value: len(self.list_by_type(t)) 
                for t in ArtifactType
            },
            "by_state": {
                s.value: len(self.list_by_state(s))
                for s in LifecycleState
            },
            "approved_count": len(self.list_approved()),
            "draft_count": len(self.list_drafts())
        }
    
    def export_state(self) -> dict:
        """Export full registry state for persistence."""
        return {
            "artifacts": {
                aid: {
                    "id": a.id,
                    "name": a.name,
                    "type": a.type.value,
                    "state": a.state.value,
                    "owner": a.owner,
                    "created_at": a.created_at,
                    "updated_at": a.updated_at,
                    "traces_to": a.traces_to,
                    "evidence_strength": a.evidence_strength,
                    "description": a.description,
                    "location": a.location
                }
                for aid, a in self.artifacts.items()
            }
        }
    
    @classmethod
    def import_state(cls, state: dict) -> "ArtifactRegistry":
        """Import registry state from persistence."""
        registry = cls()
        for aid, adata in state.get("artifacts", {}).items():
            artifact = Artifact(
                id=adata["id"],
                name=adata["name"],
                type=ArtifactType(adata["type"]),
                state=LifecycleState(adata["state"]),
                owner=adata.get("owner", ""),
                created_at=adata.get("created_at", ""),
                updated_at=adata.get("updated_at", ""),
                traces_to=adata.get("traces_to", []),
                evidence_strength=adata.get("evidence_strength", ""),
                description=adata.get("description", ""),
                location=adata.get("location", "")
            )
            registry.register(artifact)
        return registry

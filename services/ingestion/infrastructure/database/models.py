"""SQLAlchemy database models."""
from datetime import datetime
from uuid import UUID
from sqlalchemy import Column, String, DateTime, Integer, Text
from sqlalchemy.dialects.postgresql import UUID as PGUUID, JSONB
from sqlalchemy.orm import declarative_base

from domain.entities.api_specification import APISpecification
from domain.entities.ingestion_result import IngestionResult
from domain.value_objects.version import Version
from domain.value_objects.source_type import SourceType

Base = declarative_base()


class APISpecificationModel(Base):
    """SQLAlchemy model for API specifications."""
    
    __tablename__ = "api_specifications"
    
    id = Column(PGUUID(as_uuid=True), primary_key=True)
    name = Column(String(255), nullable=False)
    version = Column(String(50), nullable=False)
    source_type = Column(String(50), nullable=False)
    source_path = Column(Text, nullable=True)
    content_hash = Column(String(64), nullable=False)
    metadata = Column(JSONB, nullable=True)
    created_at = Column(DateTime, nullable=False)
    updated_at = Column(DateTime, nullable=False)
    created_by = Column(PGUUID(as_uuid=True), nullable=True)
    
    def to_entity(self) -> APISpecification:
        """Convert model to domain entity."""
        return APISpecification(
            id=self.id,
            name=self.name,
            version=Version.parse(self.version),
            source_type=SourceType(self.source_type),
            source_path=self.source_path,
            content_hash=self.content_hash,
            metadata=self.metadata or {},
            created_at=self.created_at,
            updated_at=self.updated_at,
            created_by=self.created_by
        )
    
    @classmethod
    def from_entity(cls, entity: APISpecification) -> "APISpecificationModel":
        """Create model from domain entity."""
        return cls(
            id=entity.id,
            name=entity.name,
            version=str(entity.version),
            source_type=entity.source_type.value,
            source_path=entity.source_path,
            content_hash=entity.content_hash,
            metadata=entity.metadata,
            created_at=entity.created_at,
            updated_at=entity.updated_at,
            created_by=entity.created_by
        )
    
    def update_from_entity(self, entity: APISpecification):
        """Update model from entity."""
        self.name = entity.name
        self.version = str(entity.version)
        self.metadata = entity.metadata
        self.updated_at = entity.updated_at


class IngestionLogModel(Base):
    """SQLAlchemy model for ingestion logs."""
    
    __tablename__ = "ingestion_logs"
    
    id = Column(PGUUID(as_uuid=True), primary_key=True)
    source_type = Column(String(50), nullable=False)
    source_path = Column(Text, nullable=True)
    status = Column(String(50), nullable=False)
    apis_ingested = Column(Integer, default=0)
    error_message = Column(Text, nullable=True)
    created_at = Column(DateTime, nullable=False)
    
    def to_entity(self) -> IngestionResult:
        """Convert model to domain entity."""
        return IngestionResult(
            id=self.id,
            source_type=self.source_type,
            source_path=self.source_path,
            status=self.status,
            apis_ingested=self.apis_ingested,
            error_message=self.error_message,
            created_at=self.created_at
        )
    
    @classmethod
    def from_entity(cls, entity: IngestionResult) -> "IngestionLogModel":
        """Create model from domain entity."""
        return cls(
            id=entity.id,
            source_type=entity.source_type,
            source_path=entity.source_path,
            status=entity.status,
            apis_ingested=entity.apis_ingested,
            error_message=entity.error_message,
            created_at=entity.created_at
        )


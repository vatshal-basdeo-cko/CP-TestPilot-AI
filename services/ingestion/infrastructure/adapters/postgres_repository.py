"""PostgreSQL repository implementation."""
import logging
from typing import List, Optional
from uuid import UUID
from sqlalchemy import select, desc
from sqlalchemy.ext.asyncio import AsyncSession, create_async_engine, async_sessionmaker

from ...domain.entities.api_specification import APISpecification
from ...domain.entities.ingestion_result import IngestionResult
from ...domain.repositories.ingestion_repository import IngestionRepository
from ..database.models import APISpecificationModel, IngestionLogModel

logger = logging.getLogger(__name__)


class PostgresRepository(IngestionRepository):
    """PostgreSQL implementation of ingestion repository."""
    
    def __init__(self, database_url: str):
        self.engine = create_async_engine(database_url, echo=False)
        self.SessionLocal = async_sessionmaker(
            self.engine, class_=AsyncSession, expire_on_commit=False
        )
    
    async def save_api_specification(self, api_spec: APISpecification) -> APISpecification:
        async with self.SessionLocal() as session:
            model = APISpecificationModel.from_entity(api_spec)
            session.add(model)
            await session.commit()
            await session.refresh(model)
            return model.to_entity()
    
    async def find_by_name_and_version(self, name: str, version: str) -> Optional[APISpecification]:
        async with self.SessionLocal() as session:
            result = await session.execute(
                select(APISpecificationModel).where(
                    APISpecificationModel.name == name,
                    APISpecificationModel.version == version
                )
            )
            model = result.scalar_one_or_none()
            return model.to_entity() if model else None
    
    async def find_by_content_hash(self, content_hash: str) -> Optional[APISpecification]:
        async with self.SessionLocal() as session:
            result = await session.execute(
                select(APISpecificationModel).where(
                    APISpecificationModel.content_hash == content_hash
                )
            )
            model = result.scalar_one_or_none()
            return model.to_entity() if model else None
    
    async def find_by_id(self, api_spec_id: UUID) -> Optional[APISpecification]:
        async with self.SessionLocal() as session:
            model = await session.get(APISpecificationModel, api_spec_id)
            return model.to_entity() if model else None
    
    async def list_all(self, limit: int = 100, offset: int = 0) -> List[APISpecification]:
        async with self.SessionLocal() as session:
            result = await session.execute(
                select(APISpecificationModel)
                .order_by(desc(APISpecificationModel.created_at))
                .limit(limit)
                .offset(offset)
            )
            models = result.scalars().all()
            return [model.to_entity() for model in models]
    
    async def update_api_specification(self, api_spec: APISpecification) -> APISpecification:
        async with self.SessionLocal() as session:
            model = await session.get(APISpecificationModel, api_spec.id)
            if not model:
                raise ValueError(f"API specification not found: {api_spec.id}")
            model.update_from_entity(api_spec)
            await session.commit()
            await session.refresh(model)
            return model.to_entity()
    
    async def delete_api_specification(self, api_spec_id: UUID) -> bool:
        async with self.SessionLocal() as session:
            model = await session.get(APISpecificationModel, api_spec_id)
            if not model:
                return False
            await session.delete(model)
            await session.commit()
            return True
    
    async def save_ingestion_result(self, result: IngestionResult) -> IngestionResult:
        async with self.SessionLocal() as session:
            model = IngestionLogModel.from_entity(result)
            session.add(model)
            await session.commit()
            await session.refresh(model)
            return model.to_entity()
    
    async def get_recent_ingestions(self, limit: int = 10) -> List[IngestionResult]:
        async with self.SessionLocal() as session:
            result = await session.execute(
                select(IngestionLogModel)
                .order_by(desc(IngestionLogModel.created_at))
                .limit(limit)
            )
            models = result.scalars().all()
            return [model.to_entity() for model in models]


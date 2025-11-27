"""Use case for ingesting API configurations from a folder."""
import logging
from pathlib import Path
from typing import List, Optional
from uuid import UUID

from ...domain.entities.ingestion_result import IngestionResult
from ...domain.value_objects.source_type import SourceType
from .ingest_from_file import IngestFromFileUseCase

logger = logging.getLogger(__name__)


class IngestFromFolderUseCase:
    """Use case for ingesting all API configurations from a folder."""
    
    def __init__(
        self,
        ingest_file_use_case: IngestFromFileUseCase,
    ):
        self.ingest_file_use_case = ingest_file_use_case
    
    async def execute(
        self,
        folder_path: str,
        recursive: bool = False,
        created_by: Optional[UUID] = None
    ) -> IngestionResult:
        """
        Ingest all API configuration files from a folder.
        
        Args:
            folder_path: Path to the folder containing configuration files
            recursive: Whether to search subdirectories
            created_by: UUID of the user creating these specifications
            
        Returns:
            IngestionResult with aggregate status and details
        """
        try:
            logger.info(f"Starting folder ingestion: {folder_path}")
            
            # Find all configuration files
            config_files = self._find_config_files(folder_path, recursive)
            
            if not config_files:
                logger.warning(f"No configuration files found in {folder_path}")
                return IngestionResult(
                    source_type=SourceType.FILE.value,
                    source_path=folder_path,
                    status="success",
                    apis_ingested=0,
                    error_message="No configuration files found"
                )
            
            logger.info(f"Found {len(config_files)} configuration files")
            
            # Ingest each file
            results = []
            api_ids = []
            errors = []
            
            for config_file in config_files:
                try:
                    result = await self.ingest_file_use_case.execute(
                        file_path=str(config_file),
                        created_by=created_by
                    )
                    results.append(result)
                    api_ids.extend(result.api_ids)
                    
                    if result.error_message:
                        errors.append(f"{config_file.name}: {result.error_message}")
                        
                except Exception as e:
                    logger.error(f"Error ingesting {config_file}: {str(e)}")
                    errors.append(f"{config_file.name}: {str(e)}")
            
            # Aggregate results
            total_ingested = sum(r.apis_ingested for r in results)
            
            # Determine overall status
            if not results:
                status = "failed"
            elif all(r.is_successful() for r in results):
                status = "success"
            elif any(r.is_successful() for r in results):
                status = "partial"
            else:
                status = "failed"
            
            error_message = "; ".join(errors) if errors else None
            
            result = IngestionResult(
                source_type=SourceType.FILE.value,
                source_path=folder_path,
                status=status,
                apis_ingested=total_ingested,
                api_ids=api_ids,
                error_message=error_message
            )
            
            logger.info(
                f"Folder ingestion complete: {total_ingested} APIs ingested, "
                f"{len(errors)} errors"
            )
            
            return result
            
        except Exception as e:
            logger.error(f"Error ingesting folder {folder_path}: {str(e)}", exc_info=True)
            return IngestionResult(
                source_type=SourceType.FILE.value,
                source_path=folder_path,
                status="failed",
                apis_ingested=0,
                error_message=str(e)
            )
    
    def _find_config_files(
        self, 
        folder_path: str, 
        recursive: bool
    ) -> List[Path]:
        """
        Find all configuration files in the folder.
        
        Looks for .yaml, .yml, and .json files.
        """
        folder = Path(folder_path)
        
        if not folder.exists():
            raise ValueError(f"Folder does not exist: {folder_path}")
        
        if not folder.is_dir():
            raise ValueError(f"Path is not a directory: {folder_path}")
        
        patterns = ["*.yaml", "*.yml", "*.json"]
        files = []
        
        for pattern in patterns:
            if recursive:
                files.extend(folder.rglob(pattern))
            else:
                files.extend(folder.glob(pattern))
        
        # Sort for consistent ordering
        return sorted(files)


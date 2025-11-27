"""Faker adapter for test data generation."""
import logging
from faker import Faker

logger = logging.getLogger(__name__)


class FakerAdapter:
    """Adapter for generating realistic test data."""
    
    def __init__(self):
        self.faker = Faker()
        logger.info("Initialized Faker adapter")
    
    async def generate(self, field_type: str, field_name: str = ""):
        """Generate test data based on field type."""
        try:
            if field_type == 'email':
                return self.faker.email()
            elif field_type == 'card_number':
                return self.faker.credit_card_number()
            elif field_type == 'amount':
                return round(self.faker.random.uniform(10, 1000), 2)
            elif field_type == 'name':
                return self.faker.name()
            elif field_type == 'address':
                return self.faker.address()
            elif field_type == 'phone':
                return self.faker.phone_number()
            elif field_type == 'date':
                return self.faker.date()
            elif field_type == 'uuid':
                return str(self.faker.uuid4())
            else:
                return self.faker.word()
        except Exception as e:
            logger.error(f"Error generating fake data: {str(e)}")
            return "test_value"


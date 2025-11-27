using System;
using System.Threading.Tasks;
using TestPilot.Validation.Domain.Entities;

namespace TestPilot.Validation.Domain.Repositories;

public interface IValidationRepository
{
    Task<ValidationRule?> FindRuleById(Guid id);
    Task<List<ValidationRule>> FindRulesByApiSpecId(Guid apiSpecId);
    Task<ValidationRule> CreateRule(ValidationRule rule);
    Task<ValidationRule> UpdateRule(ValidationRule rule);
    Task DeleteRule(Guid id);
    Task<List<ValidationRule>> ListRules();
    Task SaveValidationResult(ValidationResult result);
}


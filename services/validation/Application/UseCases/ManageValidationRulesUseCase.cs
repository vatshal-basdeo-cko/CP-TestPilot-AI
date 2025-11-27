using System;
using System.Threading.Tasks;
using TestPilot.Validation.Domain.Entities;
using TestPilot.Validation.Domain.Repositories;

namespace TestPilot.Validation.Application.UseCases;

public class ManageValidationRulesUseCase
{
    private readonly IValidationRepository _repository;

    public ManageValidationRulesUseCase(IValidationRepository repository)
    {
        _repository = repository;
    }

    public async Task<ValidationRule> CreateRule(ValidationRule rule)
    {
        if (!rule.IsValid())
        {
            throw new ArgumentException("Invalid validation rule");
        }

        return await _repository.CreateRule(rule);
    }

    public async Task<ValidationRule?> GetRuleById(Guid id)
    {
        return await _repository.FindRuleById(id);
    }

    public async Task<List<ValidationRule>> GetRulesByApiSpecId(Guid apiSpecId)
    {
        return await _repository.FindRulesByApiSpecId(apiSpecId);
    }

    public async Task<ValidationRule> UpdateRule(ValidationRule rule)
    {
        if (!rule.IsValid())
        {
            throw new ArgumentException("Invalid validation rule");
        }

        rule.UpdatedAt = DateTime.UtcNow;
        return await _repository.UpdateRule(rule);
    }

    public async Task DeleteRule(Guid id)
    {
        await _repository.DeleteRule(id);
    }

    public async Task<List<ValidationRule>> ListRules()
    {
        return await _repository.ListRules();
    }
}


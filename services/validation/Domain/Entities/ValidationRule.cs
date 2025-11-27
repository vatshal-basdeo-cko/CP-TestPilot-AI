using System;

namespace TestPilot.Validation.Domain.Entities;

public class ValidationRule
{
    public Guid Id { get; set; }
    public Guid ApiSpecId { get; set; }
    public string RuleType { get; set; } = string.Empty; // schema, status, custom
    public string RuleDefinition { get; set; } = string.Empty; // JSON
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }

    public ValidationRule()
    {
        Id = Guid.NewGuid();
        CreatedAt = DateTime.UtcNow;
        UpdatedAt = DateTime.UtcNow;
    }

    public bool IsValid()
    {
        return !string.IsNullOrEmpty(RuleType) && 
               !string.IsNullOrEmpty(RuleDefinition);
    }
}


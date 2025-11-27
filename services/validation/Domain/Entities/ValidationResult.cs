using System;
using System.Collections.Generic;

namespace TestPilot.Validation.Domain.Entities;

public class ValidationResult
{
    public Guid Id { get; set; }
    public Guid ExecutionId { get; set; }
    public bool IsValid { get; set; }
    public List<ValidationError> Errors { get; set; } = new();
    public Dictionary<string, object> Details { get; set; } = new();
    public DateTime ValidatedAt { get; set; }

    public ValidationResult()
    {
        Id = Guid.NewGuid();
        ValidatedAt = DateTime.UtcNow;
    }

    public void AddError(string field, string message, string severity = "error")
    {
        Errors.Add(new ValidationError
        {
            Field = field,
            Message = message,
            Severity = severity
        });
        IsValid = false;
    }
}

public class ValidationError
{
    public string Field { get; set; } = string.Empty;
    public string Message { get; set; } = string.Empty;
    public string Severity { get; set; } = "error"; // error, warning
}


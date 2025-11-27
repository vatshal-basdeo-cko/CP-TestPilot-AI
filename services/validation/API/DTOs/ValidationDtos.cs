namespace TestPilot.Validation.API.DTOs;

public class ValidateRequestDto
{
    public object Response { get; set; } = new();
    public int StatusCode { get; set; }
    public string? JsonSchema { get; set; }
    public int? ExpectedStatusCode { get; set; }
    public Guid? ExecutionId { get; set; }
}

public class ValidationRuleDto
{
    public Guid Id { get; set; }
    public Guid ApiSpecId { get; set; }
    public string RuleType { get; set; } = string.Empty;
    public string RuleDefinition { get; set; } = string.Empty;
}

public class HealthResponseDto
{
    public string Status { get; set; } = "healthy";
    public string Service { get; set; } = "validation";
    public DateTime Timestamp { get; set; } = DateTime.UtcNow;
}


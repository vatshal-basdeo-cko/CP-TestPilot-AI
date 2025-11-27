using System;
using System.Text.Json;
using System.Threading.Tasks;
using NJsonSchema;
using TestPilot.Validation.Domain.Entities;

namespace TestPilot.Validation.Application.UseCases;

public class ValidateResponseUseCase
{
    public async Task<ValidationResult> Execute(
        object response,
        int statusCode,
        string? jsonSchema = null,
        int? expectedStatusCode = null)
    {
        var result = new ValidationResult();

        try
        {
            // Validate status code
            if (expectedStatusCode.HasValue)
            {
                if (statusCode != expectedStatusCode.Value)
                {
                    result.AddError(
                        "status_code",
                        $"Expected status code {expectedStatusCode.Value}, got {statusCode}"
                    );
                }
                else
                {
                    result.Details["status_code"] = "Valid";
                }
            }

            // Validate JSON schema if provided
            if (!string.IsNullOrEmpty(jsonSchema))
            {
                var schema = await JsonSchema.FromJsonAsync(jsonSchema);
                var responseJson = JsonSerializer.Serialize(response);
                var errors = schema.Validate(responseJson);

                if (errors.Count > 0)
                {
                    foreach (var error in errors)
                    {
                        result.AddError(
                            error.Path,
                            error.ToString()
                        );
                    }
                }
                else
                {
                    result.Details["schema"] = "Valid";
                }
            }

            // If no errors, mark as valid
            if (result.Errors.Count == 0)
            {
                result.IsValid = true;
            }
        }
        catch (Exception ex)
        {
            result.AddError("validation", $"Validation error: {ex.Message}");
        }

        return result;
    }
}


using Microsoft.AspNetCore.Mvc;
using TestPilot.Validation.API.DTOs;
using TestPilot.Validation.Application.UseCases;
using TestPilot.Validation.Domain.Entities;
using TestPilot.Validation.Domain.Repositories;

namespace TestPilot.Validation.API.Controllers;

[ApiController]
[Route("api/v1")]
public class ValidationController : ControllerBase
{
    private readonly ValidateResponseUseCase _validateUseCase;
    private readonly ManageValidationRulesUseCase _rulesUseCase;
    private readonly IValidationRepository _repository;
    private readonly ILogger<ValidationController> _logger;

    public ValidationController(
        ValidateResponseUseCase validateUseCase,
        ManageValidationRulesUseCase rulesUseCase,
        IValidationRepository repository,
        ILogger<ValidationController> logger)
    {
        _validateUseCase = validateUseCase;
        _rulesUseCase = rulesUseCase;
        _repository = repository;
        _logger = logger;
    }

    [HttpPost("validate")]
    public async Task<IActionResult> ValidateResponse([FromBody] ValidateRequestDto request)
    {
        try
        {
            var result = await _validateUseCase.Execute(
                request.Response,
                request.StatusCode,
                request.JsonSchema,
                request.ExpectedStatusCode
            );

            // Save validation result if execution ID provided
            if (request.ExecutionId.HasValue)
            {
                result.ExecutionId = request.ExecutionId.Value;
                await _repository.SaveValidationResult(result);
            }

            return Ok(result);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Validation error");
            return StatusCode(500, new { error = ex.Message });
        }
    }

    [HttpGet("rules")]
    public async Task<IActionResult> ListRules()
    {
        var rules = await _rulesUseCase.ListRules();
        return Ok(new { rules, count = rules.Count });
    }

    [HttpGet("rules/{id}")]
    public async Task<IActionResult> GetRule(Guid id)
    {
        var rule = await _rulesUseCase.GetRuleById(id);
        if (rule == null)
        {
            return NotFound(new { error = "Rule not found" });
        }
        return Ok(rule);
    }

    [HttpPost("rules")]
    public async Task<IActionResult> CreateRule([FromBody] ValidationRuleDto dto)
    {
        try
        {
            var rule = new ValidationRule
            {
                ApiSpecId = dto.ApiSpecId,
                RuleType = dto.RuleType,
                RuleDefinition = dto.RuleDefinition
            };

            var created = await _rulesUseCase.CreateRule(rule);
            return CreatedAtAction(nameof(GetRule), new { id = created.Id }, created);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error creating rule");
            return BadRequest(new { error = ex.Message });
        }
    }

    [HttpPut("rules/{id}")]
    public async Task<IActionResult> UpdateRule(Guid id, [FromBody] ValidationRuleDto dto)
    {
        try
        {
            var rule = await _rulesUseCase.GetRuleById(id);
            if (rule == null)
            {
                return NotFound(new { error = "Rule not found" });
            }

            rule.RuleType = dto.RuleType;
            rule.RuleDefinition = dto.RuleDefinition;

            var updated = await _rulesUseCase.UpdateRule(rule);
            return Ok(updated);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error updating rule");
            return BadRequest(new { error = ex.Message });
        }
    }

    [HttpDelete("rules/{id}")]
    public async Task<IActionResult> DeleteRule(Guid id)
    {
        try
        {
            await _rulesUseCase.DeleteRule(id);
            return Ok(new { message = "Rule deleted successfully" });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error deleting rule");
            return StatusCode(500, new { error = ex.Message });
        }
    }

    [HttpGet("/health")]
    public IActionResult HealthCheck()
    {
        return Ok(new HealthResponseDto());
    }
}


using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using Npgsql;
using TestPilot.Validation.Domain.Entities;
using TestPilot.Validation.Domain.Repositories;
using System.Text.Json;

namespace TestPilot.Validation.Infrastructure.Adapters;

public class PostgresRepository : IValidationRepository
{
    private readonly string _connectionString;

    public PostgresRepository(string connectionString)
    {
        _connectionString = connectionString;
    }

    public async Task<ValidationRule> CreateRule(ValidationRule rule)
    {
        using var conn = new NpgsqlConnection(_connectionString);
        await conn.OpenAsync();

        var sql = @"
            INSERT INTO validation_rules (id, api_spec_id, rule_type, rule_definition, created_at, updated_at)
            VALUES (@id, @api_spec_id, @rule_type, @rule_definition, @created_at, @updated_at)
            RETURNING *";

        using var cmd = new NpgsqlCommand(sql, conn);
        cmd.Parameters.AddWithValue("id", rule.Id);
        cmd.Parameters.AddWithValue("api_spec_id", rule.ApiSpecId);
        cmd.Parameters.AddWithValue("rule_type", rule.RuleType);
        cmd.Parameters.AddWithValue("rule_definition", rule.RuleDefinition);
        cmd.Parameters.AddWithValue("created_at", rule.CreatedAt);
        cmd.Parameters.AddWithValue("updated_at", rule.UpdatedAt);

        await cmd.ExecuteNonQueryAsync();
        return rule;
    }

    public async Task<ValidationRule?> FindRuleById(Guid id)
    {
        using var conn = new NpgsqlConnection(_connectionString);
        await conn.OpenAsync();

        var sql = "SELECT * FROM validation_rules WHERE id = @id";
        using var cmd = new NpgsqlCommand(sql, conn);
        cmd.Parameters.AddWithValue("id", id);

        using var reader = await cmd.ExecuteReaderAsync();
        if (await reader.ReadAsync())
        {
            return MapToValidationRule(reader);
        }

        return null;
    }

    public async Task<List<ValidationRule>> FindRulesByApiSpecId(Guid apiSpecId)
    {
        using var conn = new NpgsqlConnection(_connectionString);
        await conn.OpenAsync();

        var sql = "SELECT * FROM validation_rules WHERE api_spec_id = @api_spec_id";
        using var cmd = new NpgsqlCommand(sql, conn);
        cmd.Parameters.AddWithValue("api_spec_id", apiSpecId);

        var rules = new List<ValidationRule>();
        using var reader = await cmd.ExecuteReaderAsync();
        while (await reader.ReadAsync())
        {
            rules.Add(MapToValidationRule(reader));
        }

        return rules;
    }

    public async Task<ValidationRule> UpdateRule(ValidationRule rule)
    {
        using var conn = new NpgsqlConnection(_connectionString);
        await conn.OpenAsync();

        var sql = @"
            UPDATE validation_rules 
            SET rule_type = @rule_type, rule_definition = @rule_definition, updated_at = @updated_at
            WHERE id = @id";

        using var cmd = new NpgsqlCommand(sql, conn);
        cmd.Parameters.AddWithValue("id", rule.Id);
        cmd.Parameters.AddWithValue("rule_type", rule.RuleType);
        cmd.Parameters.AddWithValue("rule_definition", rule.RuleDefinition);
        cmd.Parameters.AddWithValue("updated_at", rule.UpdatedAt);

        await cmd.ExecuteNonQueryAsync();
        return rule;
    }

    public async Task DeleteRule(Guid id)
    {
        using var conn = new NpgsqlConnection(_connectionString);
        await conn.OpenAsync();

        var sql = "DELETE FROM validation_rules WHERE id = @id";
        using var cmd = new NpgsqlCommand(sql, conn);
        cmd.Parameters.AddWithValue("id", id);

        await cmd.ExecuteNonQueryAsync();
    }

    public async Task<List<ValidationRule>> ListRules()
    {
        using var conn = new NpgsqlConnection(_connectionString);
        await conn.OpenAsync();

        var sql = "SELECT * FROM validation_rules ORDER BY created_at DESC";
        using var cmd = new NpgsqlCommand(sql, conn);

        var rules = new List<ValidationRule>();
        using var reader = await cmd.ExecuteReaderAsync();
        while (await reader.ReadAsync())
        {
            rules.Add(MapToValidationRule(reader));
        }

        return rules;
    }

    public async Task SaveValidationResult(ValidationResult result)
    {
        using var conn = new NpgsqlConnection(_connectionString);
        await conn.OpenAsync();

        var sql = @"
            UPDATE test_executions 
            SET validation_result = @validation_result::jsonb
            WHERE id = @execution_id";

        using var cmd = new NpgsqlCommand(sql, conn);
        cmd.Parameters.AddWithValue("execution_id", result.ExecutionId);
        cmd.Parameters.AddWithValue("validation_result", JsonSerializer.Serialize(result));

        await cmd.ExecuteNonQueryAsync();
    }

    private ValidationRule MapToValidationRule(NpgsqlDataReader reader)
    {
        return new ValidationRule
        {
            Id = reader.GetGuid(0),
            ApiSpecId = reader.GetGuid(1),
            RuleType = reader.GetString(2),
            RuleDefinition = reader.GetString(3),
            CreatedAt = reader.GetDateTime(4),
            UpdatedAt = reader.GetDateTime(5)
        };
    }
}


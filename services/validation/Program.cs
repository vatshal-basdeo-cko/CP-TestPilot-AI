using Serilog;
using TestPilot.Validation.Application.UseCases;
using TestPilot.Validation.Domain.Repositories;
using TestPilot.Validation.Infrastructure.Adapters;

var builder = WebApplication.CreateBuilder(args);

// Configure Serilog
Log.Logger = new LoggerConfiguration()
    .WriteTo.Console()
    .CreateLogger();

builder.Host.UseSerilog();

// Add services
builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

// CORS
builder.Services.AddCors(options =>
{
    options.AddDefaultPolicy(policy =>
    {
        policy.AllowAnyOrigin()
              .AllowAnyMethod()
              .AllowAnyHeader();
    });
});

// Register dependencies
var connectionString = builder.Configuration.GetConnectionString("DefaultConnection") 
    ?? "Host=postgres;Database=testpilot;Username=testpilot;Password=testpilot";

builder.Services.AddSingleton<IValidationRepository>(
    new PostgresRepository(connectionString));
builder.Services.AddScoped<ValidateResponseUseCase>();
builder.Services.AddScoped<ManageValidationRulesUseCase>();

var app = builder.Build();

// Configure middleware
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseCors();
app.UseAuthorization();
app.MapControllers();

Log.Information("Starting Validation Service on port 8004");

app.Run();


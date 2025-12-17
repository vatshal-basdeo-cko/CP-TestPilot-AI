# Frontend Placeholder - Will be replaced with React app
FROM nginx:alpine

# Create a simple placeholder page
RUN echo '<!DOCTYPE html>\
<html lang="en">\
<head>\
    <meta charset="UTF-8">\
    <meta name="viewport" content="width=device-width, initial-scale=1.0">\
    <title>TestPilot AI</title>\
    <style>\
        * { margin: 0; padding: 0; box-sizing: border-box; }\
        body {\
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;\
            background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);\
            min-height: 100vh;\
            display: flex;\
            align-items: center;\
            justify-content: center;\
            color: white;\
        }\
        .container {\
            text-align: center;\
            padding: 2rem;\
        }\
        h1 {\
            font-size: 3rem;\
            margin-bottom: 1rem;\
            background: linear-gradient(90deg, #00d9ff, #00ff88);\
            -webkit-background-clip: text;\
            -webkit-text-fill-color: transparent;\
        }\
        .subtitle {\
            font-size: 1.2rem;\
            color: #a0a0a0;\
            margin-bottom: 2rem;\
        }\
        .status {\
            background: rgba(255,255,255,0.1);\
            padding: 1.5rem 2rem;\
            border-radius: 12px;\
            margin-bottom: 2rem;\
        }\
        .status h2 {\
            color: #00ff88;\
            margin-bottom: 0.5rem;\
        }\
        .api-info {\
            background: rgba(0,0,0,0.3);\
            padding: 1rem;\
            border-radius: 8px;\
            font-family: monospace;\
            margin-top: 1rem;\
        }\
        .api-info code {\
            color: #00d9ff;\
        }\
        a {\
            color: #00d9ff;\
            text-decoration: none;\
        }\
        a:hover {\
            text-decoration: underline;\
        }\
    </style>\
</head>\
<body>\
    <div class="container">\
        <h1>🚀 TestPilot AI</h1>\
        <p class="subtitle">AI-Powered API Testing Platform</p>\
        <div class="status">\
            <h2>✅ Backend Services Running</h2>\
            <p>Frontend UI coming soon!</p>\
        </div>\
        <div class="api-info">\
            <p>API Gateway: <code><a href="http://localhost:8000/health">http://localhost:8000</a></code></p>\
            <p style="margin-top: 0.5rem;">Health Check: <code><a href="http://localhost:8000/health/all">/health/all</a></code></p>\
        </div>\
    </div>\
</body>\
</html>' > /usr/share/nginx/html/index.html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]


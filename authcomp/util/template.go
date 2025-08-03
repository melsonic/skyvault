package util

const EmailTemplate = `
<h3>Password Reset<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Password Reset</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            max-width: 500px;
            margin: 0 auto;
            padding: 20px;
            color: #333;
        }
        h1 {
            font-size: 24px;
            margin-bottom: 20px;
        }
        .divider {
            border-top: 1px solid #e0e0e0;
            margin: 20px 0;
        }
        .button {
            display: inline-block;
            padding: 10px 20px;
            background-color: #007bff;
            color: white;
            text-decoration: none;
            border-radius: 4px;
            margin: 10px 0;
        }
        .footer {
            font-size: 14px;
            color: #666;
            margin-top: 20px;
        }
    </style>
</head>
<body>
    <h1>Password Reset</h1>
    <p>Hi,<br>
    We got a request to reset your SkyVault password.</p>
    
    <div class="divider"></div>
    
    <h2>Reset password</h2>
    <a href="{{.ResetURL}}" class="button">Reset Password</a>
    
    <p>If you ignore this message, your <strong>password</strong> won't be changed.</p>
    <p>If you didn't request a <strong>password reset</strong>, let us know.</p>
    
    <div class="divider"></div>
    
    <p class="footer">If the button above does not appear, please copy and paste this link into your browser's address bar:<br>
    {{.ResetURL}}</p>
</body>
</html>set</h3>
`
const express = require('express');
const mysql = require('mysql');
const { exec } = require('child_process');
const fs = require('fs');
const app = express();

// Hardcoded credentials
const dbConfig = {
    host: 'localhost',
    user: 'root',
    password: 'MySecretPassword123!',  // VULNERABLE: Hardcoded password
    database: 'myapp'
};

const connection = mysql.createConnection(dbConfig);

// SQL Injection vulnerability
app.get('/user', (req, res) => {
    const userId = req.query.id;

    // VULNERABLE: String concatenation in SQL query
    const query = `SELECT * FROM users WHERE id = ${userId}`;

    connection.query(query, (error, results) => {
        if (error) throw error;
        res.json(results);
    });
});

// NoSQL Injection vulnerability
app.get('/search', (req, res) => {
    const username = req.query.username;

    // VULNERABLE: User input directly in MongoDB query
    const query = { username: username };
    db.collection('users').find(query).toArray((err, docs) => {
        res.json(docs);
    });
});

// Command Injection vulnerability
app.get('/ping', (req, res) => {
    const host = req.query.host;

    // VULNERABLE: User input in shell command
    exec(`ping -c 4 ${host}`, (error, stdout, stderr) => {
        res.send(stdout);
    });
});

// Path Traversal vulnerability
app.get('/file', (req, res) => {
    const filename = req.query.name;

    // VULNERABLE: No path validation
    fs.readFile(`/var/www/uploads/${filename}`, 'utf8', (err, data) => {
        if (err) {
            res.status(500).send('Error reading file');
            return;
        }
        res.send(data);
    });
});

// XSS vulnerability
app.get('/display', (req, res) => {
    const message = req.query.msg;

    // VULNERABLE: Unescaped user input in HTML response
    res.send(`<html><body><h1>${message}</h1></body></html>`);
});

// Insecure Direct Object Reference (IDOR)
app.get('/document', (req, res) => {
    const docId = req.query.id;

    // VULNERABLE: No authorization check
    const query = `SELECT * FROM documents WHERE id = ${docId}`;
    connection.query(query, (error, results) => {
        res.json(results);
    });
});

// Weak cryptography
const crypto = require('crypto');

function hashPassword(password) {
    // VULNERABLE: Using weak MD5 algorithm
    return crypto.createHash('md5').update(password).digest('hex');
}

// Insecure random number generation
function generateToken() {
    // VULNERABLE: Math.random() is not cryptographically secure
    return Math.random().toString(36).substring(7);
}

// Sensitive data exposure
app.get('/config', (req, res) => {
    // VULNERABLE: Exposing sensitive configuration
    res.json({
        database: dbConfig,
        apiKey: '123',
        awsSecretKey: '456'
    });
});

// SSRF vulnerability
const axios = require('axios');

app.get('/fetch', async (req, res) => {
    const url = req.query.url;

    // VULNERABLE: Fetching user-provided URL
    try {
        const response = await axios.get(url);
        res.send(response.data);
    } catch (error) {
        res.status(500).send('Error fetching URL');
    }
});

// Insecure cookie settings
app.get('/login', (req, res) => {
    // VULNERABLE: Cookie without secure flags
    res.cookie('session', '12345', {
        httpOnly: false,  // Should be true
        secure: false,    // Should be true
        sameSite: 'none'  // Vulnerable to CSRF
    });
    res.send('Logged in');
});

app.listen(3000, () => {
    console.log('Server running on port 3000');
});

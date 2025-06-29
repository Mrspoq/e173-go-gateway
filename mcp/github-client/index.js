const express = require('express');
const bodyParser = require('body-parser');
const { Octokit } = require('@octokit/rest');
const fs = require('fs');
const path = require('path');

const app = express();
app.use(bodyParser.json());

// Load GitHub configuration
const configPath = path.join(__dirname, '../../.github_config');
const config = fs.readFileSync(configPath, 'utf8')
  .split('\n')
  .reduce((acc, line) => {
    const [key, value] = line.split('=');
    if (key && value) {
      acc[key.replace('GITHUB_', '').toLowerCase()] = value;
    }
    return acc;
  }, {});

const octokit = new Octokit({
  auth: config.token
});

// Create project
app.post('/project/create', async (req, res) => {
  try {
    const { name, description } = req.body;
    const result = await octokit.projects.createForUser({
      name,
      body: description
    });
    res.json({ success: true, data: result.data });
  } catch (error) {
    res.status(500).json({ success: false, error: error.message });
  }
});

// Create milestone
app.post('/milestone/create', async (req, res) => {
  try {
    const { title, description, due_on } = req.body;
    const result = await octokit.issues.createMilestone({
      owner: config.user,
      repo: config.repo,
      title,
      description,
      due_on
    });
    res.json({ success: true, data: result.data });
  } catch (error) {
    res.status(500).json({ success: false, error: error.message });
  }
});

// Create issue
app.post('/issue/create', async (req, res) => {
  try {
    const { title, body, labels, milestone } = req.body;
    const result = await octokit.issues.create({
      owner: config.user,
      repo: config.repo,
      title,
      body,
      labels: labels ? labels.split(',') : [],
      milestone
    });
    res.json({ success: true, data: result.data });
  } catch (error) {
    res.status(500).json({ success: false, error: error.message });
  }
});

// List repositories
app.get('/repos', async (req, res) => {
  try {
    const result = await octokit.repos.listForAuthenticatedUser();
    res.json({ success: true, data: result.data });
  } catch (error) {
    res.status(500).json({ success: false, error: error.message });
  }
});

const PORT = process.env.PORT || 3002;
app.listen(PORT, () => {
  console.log(`GitHub MCP Server running on port ${PORT}`);
  console.log(`Using GitHub user: ${config.user}`);
  console.log(`Using repository: ${config.repo}`);
});
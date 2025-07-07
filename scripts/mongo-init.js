// MongoDB initialization script
// This script runs automatically when the MongoDB container starts for the first time

// Switch to the odds_worklog_db database
db = db.getSiblingDB("odds_worklog_db");

// Create admin user for odds_worklog_db database
db.createUser({
  user: "admin",
  pwd: "admin",
  roles: [{ role: "readWrite", db: "odds_worklog_db" }],
});

print("Admin user created successfully for odds_worklog_db database");

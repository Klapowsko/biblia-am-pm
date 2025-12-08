-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create reading_plans table
CREATE TABLE IF NOT EXISTS reading_plans (
    id SERIAL PRIMARY KEY,
    day_of_year INTEGER NOT NULL UNIQUE,
    old_testament_ref VARCHAR(255),
    new_testament_ref VARCHAR(255),
    psalms_ref VARCHAR(255),
    proverbs_ref VARCHAR(255)
);

-- Create user_progress table
CREATE TABLE IF NOT EXISTS user_progress (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reading_plan_id INTEGER NOT NULL REFERENCES reading_plans(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    morning_completed BOOLEAN DEFAULT FALSE,
    evening_completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP,
    UNIQUE(user_id, reading_plan_id, date)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_user_progress_user_id ON user_progress(user_id);
CREATE INDEX IF NOT EXISTS idx_user_progress_date ON user_progress(date);
CREATE INDEX IF NOT EXISTS idx_reading_plans_day_of_year ON reading_plans(day_of_year);


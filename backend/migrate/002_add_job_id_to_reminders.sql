ALTER TABLE reminders ADD COLUMN job_id BIGINT;

CREATE INDEX idx_reminders_job_id ON reminders(job_id);
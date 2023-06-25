CREATE TABLE IF EXISTS monitors (
    monitor_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_email TEXT NOT NULL,
    type TEXT NOT NULL,
    url TEXT NOT NULL,
    method TEXT NOT NULL,
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    body TEXT ,
    headers TEXT,
    parameters TEXT,
    description TEXT,
    frequency_minutes INTEGER,
    threshold_minutes INTEGER 
    CONSTRAINT monitor_keys PRIMARY (monitor_id, user_email, type, url, method)
);
```


-- type Monitor struct {
-- 	MonitorID        int64     `json:"monitor_id" `
-- 	UserEmail        string    `json:"user_email"`
-- 	MonitorType      string    `json:"type"`
-- 	URL              string    `json:"url"`
-- 	Method           string    `json:"method"`
-- 	UpdatedAt        time.Time `json:"updated_at"`
-- 	Body             string    `json:"body"`
-- 	Headers          string    `json:"headers"`
-- 	Parameters       string    `json:"parameters"`
-- 	Description      string    `json:"description"`
-- 	FrequencyMinutes int       `json:"frequency_minutes"`
-- 	ThresholdMinutes int       `json:"threshold_minutes"`
-- }
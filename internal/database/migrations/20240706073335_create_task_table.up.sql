CREATE TABLE IF NOT EXISTS Task(
   id SERIAL PRIMARY KEY,
   people_id integer
       REFERENCES People (id)
           ON DELETE CASCADE
           ON UPDATE CASCADE,
   task char(255),
   task_start timestamp,
   task_end timestamp NULL,
   task_interval interval NULL
)
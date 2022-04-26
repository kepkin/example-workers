CREATE TABLE targets (
	id SERIAL PRIMARY KEY,
	name varchar UNIQUE NOT NULL
);

CREATE TABLE monitor_tasks (
	id uuid PRIMARY KEY,
	start timestamp NOT NULL,
	stop  timestamp NOT NULL,
	target_id bigint NOT NULL,
	frequency int NOT NULL,
	status    int NOT NULL,
	liveprobe timestamp,
	CONSTRAINT fk_target FOREIGN KEY (target_id) REFERENCES targets(id)
	-- TODO 
	-- index on start stop
);

CREATE TABLE monitor_results (
	monitor_task_id uuid,
	time timestamp  NOT NULL,
	price decimal(10, 2) NOT NULL,
	CONSTRAINT fk_monitor_task FOREIGN KEY (monitor_task_id) REFERENCES monitor_tasks(id),
	PRIMARY KEY(monitor_task_id, time)
);


INSERT INTO targets (name) VALUES('bitcoin');
CREATE TABLE public.subscriptions (
	user_id uuid NOT NULL,
	service_name varchar NOT NULL,
	price integer NOT NULL,
	start_date date NOT NULL,
	end_date date NULL,
	CONSTRAINT subscriptions_pk PRIMARY KEY (user_id,service_name)
);
CREATE TABLE users (
   id uuid NOT NULL,
   login character varying(45) NOT NULL,
   password character varying(45),
   created_at timestamp without time zone NOT NULL DEFAULT now(),
   updated_at timestamp without time zone NOT NULL DEFAULT now(),
   last_login timestamp without time zone,
   PRIMARY KEY (id),
   CONSTRAINT login_unq UNIQUE (login)
)
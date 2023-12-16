CREATE TABLE public.meteodata (
	id bigserial NOT NULL,
	dt timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    temperature decimal,
    pressure decimal,
    altitude decimal
);
CREATE INDEX meteodata_datetime_idx ON public.meteodata USING btree (dt);

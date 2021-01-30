--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.17
-- Dumped by pg_dump version 12.3 (Ubuntu 12.3-1.pgdg16.04+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

--
-- Name: coil; Type: TABLE; Schema: public; Owner: electrolab
--

CREATE TABLE public.coil (
    id integer NOT NULL,
    transformer integer NOT NULL,
    coilnumber integer NOT NULL,
    tap integer NOT NULL,
    coiltype integer,
    classaccuracy character varying(10) NOT NULL,
    primarycurrent numeric(8,2),
    secondcurrent numeric(8,2),
    secondload numeric(8,2) NOT NULL,
    magneticvoltage numeric(8,2),
    magneticcurrent numeric(8,2),
    resistance numeric(8,2),
    rating character varying(10),
    quadroload numeric(8,2),
    ampereturn integer
);


ALTER TABLE public.coil OWNER TO electrolab;

--
-- Name: coil_id_seq; Type: SEQUENCE; Schema: public; Owner: electrolab
--

CREATE SEQUENCE public.coil_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.coil_id_seq OWNER TO electrolab;

--
-- Name: coil_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: electrolab
--

ALTER SEQUENCE public.coil_id_seq OWNED BY public.coil.id;


--
-- Name: serial_number; Type: TABLE; Schema: public; Owner: electrolab
--

CREATE TABLE public.serial_number (
    id integer NOT NULL,
    ordernumber character varying(20) NOT NULL,
    series character varying(20) NOT NULL,
    serialnumber integer NOT NULL,
    makedate integer NOT NULL,
    transformer integer NOT NULL,
    replace boolean
);


ALTER TABLE public.serial_number OWNER TO electrolab;

--
-- Name: COLUMN serial_number.replace; Type: COMMENT; Schema: public; Owner: electrolab
--

COMMENT ON COLUMN public.serial_number.replace IS 'Трансформатор изготавливается взамен бракованому';


--
-- Name: serial_number_id_seq; Type: SEQUENCE; Schema: public; Owner: electrolab
--

CREATE SEQUENCE public.serial_number_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.serial_number_id_seq OWNER TO electrolab;

--
-- Name: serial_number_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: electrolab
--

ALTER SEQUENCE public.serial_number_id_seq OWNED BY public.serial_number.id;


--
-- Name: transformer; Type: TABLE; Schema: public; Owner: electrolab
--

CREATE TABLE public.transformer (
    id integer NOT NULL,
    fullname character varying(150) NOT NULL,
    shortname character varying(150) NOT NULL,
    type character varying(10) NOT NULL,
    standart character varying(100),
    voltage numeric(6,2),
    maxopervoltage numeric(8,2),
    frequency integer,
    quantsecondcoil integer,
    isolationlevel character(1),
    climat character varying(3),
    weight numeric(6,2)
);


ALTER TABLE public.transformer OWNER TO electrolab;

--
-- Name: transformer_id_seq; Type: SEQUENCE; Schema: public; Owner: electrolab
--

CREATE SEQUENCE public.transformer_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.transformer_id_seq OWNER TO electrolab;

--
-- Name: transformer_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: electrolab
--

ALTER SEQUENCE public.transformer_id_seq OWNED BY public.transformer.id;


--
-- Name: coil id; Type: DEFAULT; Schema: public; Owner: electrolab
--

ALTER TABLE ONLY public.coil ALTER COLUMN id SET DEFAULT nextval('public.coil_id_seq'::regclass);


--
-- Name: serial_number id; Type: DEFAULT; Schema: public; Owner: electrolab
--

ALTER TABLE ONLY public.serial_number ALTER COLUMN id SET DEFAULT nextval('public.serial_number_id_seq'::regclass);


--
-- Name: transformer id; Type: DEFAULT; Schema: public; Owner: electrolab
--

ALTER TABLE ONLY public.transformer ALTER COLUMN id SET DEFAULT nextval('public.transformer_id_seq'::regclass);


--
-- Name: coil coil_pkey; Type: CONSTRAINT; Schema: public; Owner: electrolab
--

ALTER TABLE ONLY public.coil
    ADD CONSTRAINT coil_pkey PRIMARY KEY (id);


--
-- Name: serial_number serial_number_pkey; Type: CONSTRAINT; Schema: public; Owner: electrolab
--

ALTER TABLE ONLY public.serial_number
    ADD CONSTRAINT serial_number_pkey PRIMARY KEY (id);


--
-- Name: serial_number serial_number_serialnumber_key; Type: CONSTRAINT; Schema: public; Owner: electrolab
--

ALTER TABLE ONLY public.serial_number
    ADD CONSTRAINT serial_number_serialnumber_key UNIQUE (serialnumber, makedate);


--
-- Name: transformer transformer_fullname_key; Type: CONSTRAINT; Schema: public; Owner: electrolab
--

ALTER TABLE ONLY public.transformer
    ADD CONSTRAINT transformer_fullname_key UNIQUE (fullname);


--
-- Name: transformer transformer_pkey; Type: CONSTRAINT; Schema: public; Owner: electrolab
--

ALTER TABLE ONLY public.transformer
    ADD CONSTRAINT transformer_pkey PRIMARY KEY (id);


--
-- Name: mdate_snum; Type: INDEX; Schema: public; Owner: electrolab
--

CREATE INDEX mdate_snum ON public.serial_number USING btree (makedate, serialnumber);


--
-- Name: coil coil_transformer_fkey; Type: FK CONSTRAINT; Schema: public; Owner: electrolab
--

ALTER TABLE ONLY public.coil
    ADD CONSTRAINT coil_transformer_fkey FOREIGN KEY (transformer) REFERENCES public.transformer(id);


--
-- Name: serial_number serial_number_transformer_fkey; Type: FK CONSTRAINT; Schema: public; Owner: electrolab
--

ALTER TABLE ONLY public.serial_number
    ADD CONSTRAINT serial_number_transformer_fkey FOREIGN KEY (transformer) REFERENCES public.transformer(id);


--
-- PostgreSQL database dump complete
--


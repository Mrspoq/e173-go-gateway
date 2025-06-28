--
-- PostgreSQL database dump
--

-- Dumped from database version 14.18 (Ubuntu 14.18-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.18 (Ubuntu 14.18-0ubuntu0.22.04.1)

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

--
-- Name: trigger_set_timestamp(); Type: FUNCTION; Schema: public; Owner: e173_user
--

CREATE FUNCTION public.trigger_set_timestamp() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$;


ALTER FUNCTION public.trigger_set_timestamp() OWNER TO e173_user;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: ai_voice_agents; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ai_voice_agents (
    id integer NOT NULL,
    agent_name character varying(100) NOT NULL,
    endpoint_url character varying(255) NOT NULL,
    api_key character varying(255),
    model_type character varying(50),
    voice_id character varying(100),
    language character varying(10) DEFAULT 'en'::character varying,
    concurrent_limit integer DEFAULT 10,
    current_sessions integer DEFAULT 0,
    cost_per_minute numeric(8,4),
    quality_rating numeric(3,2),
    active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.ai_voice_agents OWNER TO postgres;

--
-- Name: ai_voice_agents_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ai_voice_agents_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ai_voice_agents_id_seq OWNER TO postgres;

--
-- Name: ai_voice_agents_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ai_voice_agents_id_seq OWNED BY public.ai_voice_agents.id;


--
-- Name: ai_voice_sessions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ai_voice_sessions (
    id integer NOT NULL,
    session_id character varying(255) NOT NULL,
    agent_id integer,
    sip_call_id integer,
    caller_number character varying(50),
    session_duration integer DEFAULT 0,
    tokens_used integer DEFAULT 0,
    revenue_generated numeric(10,4) DEFAULT 0,
    conversation_log jsonb,
    started_at timestamp without time zone DEFAULT now(),
    ended_at timestamp without time zone
);


ALTER TABLE public.ai_voice_sessions OWNER TO postgres;

--
-- Name: ai_voice_sessions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ai_voice_sessions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.ai_voice_sessions_id_seq OWNER TO postgres;

--
-- Name: ai_voice_sessions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ai_voice_sessions_id_seq OWNED BY public.ai_voice_sessions.id;


--
-- Name: audit_logs; Type: TABLE; Schema: public; Owner: e173_user
--

CREATE TABLE public.audit_logs (
    id bigint NOT NULL,
    user_id bigint,
    action character varying(100) NOT NULL,
    entity_type character varying(50),
    entity_id bigint,
    old_values jsonb,
    new_values jsonb,
    ip_address inet,
    user_agent text,
    success boolean DEFAULT true,
    error_message text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.audit_logs OWNER TO e173_user;

--
-- Name: audit_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: e173_user
--

CREATE SEQUENCE public.audit_logs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.audit_logs_id_seq OWNER TO e173_user;

--
-- Name: audit_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: e173_user
--

ALTER SEQUENCE public.audit_logs_id_seq OWNED BY public.audit_logs.id;


--
-- Name: blacklist; Type: TABLE; Schema: public; Owner: e173_user
--

CREATE TABLE public.blacklist (
    id bigint NOT NULL,
    number_pattern character varying(50) NOT NULL,
    blacklist_type character varying(20) DEFAULT 'number'::character varying,
    reason character varying(255),
    auto_added boolean DEFAULT false,
    detection_method character varying(50),
    block_inbound boolean DEFAULT true,
    block_outbound boolean DEFAULT false,
    temporary_until timestamp with time zone,
    violation_count integer DEFAULT 1,
    last_violation_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    created_by bigint,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    reason_code character varying(50),
    expires_at timestamp without time zone,
    confidence_score numeric(3,2),
    CONSTRAINT blacklist_blacklist_type_check CHECK (((blacklist_type)::text = ANY ((ARRAY['number'::character varying, 'pattern'::character varying, 'prefix'::character varying])::text[])))
);


ALTER TABLE public.blacklist OWNER TO e173_user;

--
-- Name: blacklist_id_seq; Type: SEQUENCE; Schema: public; Owner: e173_user
--

CREATE SEQUENCE public.blacklist_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.blacklist_id_seq OWNER TO e173_user;

--
-- Name: blacklist_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: e173_user
--

ALTER SEQUENCE public.blacklist_id_seq OWNED BY public.blacklist.id;


--
-- Name: call_detail_records; Type: TABLE; Schema: public; Owner: e173_user
--

CREATE TABLE public.call_detail_records (
    id bigint NOT NULL,
    sim_card_id bigint,
    modem_id bigint,
    asterisk_unique_id character varying(128) NOT NULL,
    call_direction character varying(10) NOT NULL,
    source_number character varying(50),
    destination_number character varying(50) NOT NULL,
    call_start_time timestamp with time zone NOT NULL,
    call_answer_time timestamp with time zone,
    call_end_time timestamp with time zone NOT NULL,
    duration_seconds integer,
    billable_duration_seconds integer,
    disposition character varying(50) NOT NULL,
    hangup_cause character varying(100),
    cost_per_minute numeric(10,4),
    total_cost numeric(10,4),
    customer_id bigint,
    is_spam boolean DEFAULT false,
    spam_reason character varying(255),
    recorded_audio_path character varying(255),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT call_detail_records_call_direction_check CHECK (((call_direction)::text = ANY ((ARRAY['inbound'::character varying, 'outbound'::character varying])::text[])))
);


ALTER TABLE public.call_detail_records OWNER TO e173_user;

--
-- Name: call_detail_records_id_seq; Type: SEQUENCE; Schema: public; Owner: e173_user
--

CREATE SEQUENCE public.call_detail_records_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.call_detail_records_id_seq OWNER TO e173_user;

--
-- Name: call_detail_records_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: e173_user
--

ALTER SEQUENCE public.call_detail_records_id_seq OWNED BY public.call_detail_records.id;


--
-- Name: call_patterns; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.call_patterns (
    id integer NOT NULL,
    phone_number character varying(50) NOT NULL,
    total_calls integer DEFAULT 0,
    answered_calls integer DEFAULT 0,
    short_calls integer DEFAULT 0,
    avg_call_duration numeric(8,2) DEFAULT 0,
    spam_score numeric(3,2) DEFAULT 0,
    last_call_time timestamp without time zone,
    pattern_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE public.call_patterns OWNER TO postgres;

--
-- Name: call_patterns_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.call_patterns_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.call_patterns_id_seq OWNER TO postgres;

--
-- Name: call_patterns_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.call_patterns_id_seq OWNED BY public.call_patterns.id;


--
-- Name: customer_rate_plans; Type: TABLE; Schema: public; Owner: e173_user
--

CREATE TABLE public.customer_rate_plans (
    id bigint NOT NULL,
    customer_id bigint NOT NULL,
    rate_plan_id bigint NOT NULL,
    effective_from timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    effective_until timestamp with time zone,
    is_active boolean DEFAULT true,
    created_by bigint,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.customer_rate_plans OWNER TO e173_user;

--
-- Name: customer_rate_plans_id_seq; Type: SEQUENCE; Schema: public; Owner: e173_user
--

CREATE SEQUENCE public.customer_rate_plans_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.customer_rate_plans_id_seq OWNER TO e173_user;

--
-- Name: customer_rate_plans_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: e173_user
--

ALTER SEQUENCE public.customer_rate_plans_id_seq OWNED BY public.customer_rate_plans.id;


--
-- Name: customers; Type: TABLE; Schema: public; Owner: e173_user
--

CREATE TABLE public.customers (
    id bigint NOT NULL,
    customer_code character varying(20) NOT NULL,
    company_name character varying(255),
    contact_person character varying(255),
    email character varying(255),
    phone character varying(50),
    address text,
    city character varying(100),
    state character varying(100),
    country character varying(100),
    postal_code character varying(20),
    billing_address text,
    billing_city character varying(100),
    billing_state character varying(100),
    billing_country character varying(100),
    billing_postal_code character varying(20),
    account_status character varying(20) DEFAULT 'active'::character varying,
    credit_limit numeric(12,4) DEFAULT 0.0,
    current_balance numeric(12,4) DEFAULT 0.0,
    monthly_limit numeric(12,4),
    timezone character varying(50) DEFAULT 'UTC'::character varying,
    preferred_currency character varying(3) DEFAULT 'USD'::character varying,
    auto_recharge_enabled boolean DEFAULT false,
    auto_recharge_threshold numeric(12,4),
    auto_recharge_amount numeric(12,4),
    notes text,
    created_by bigint,
    assigned_to bigint,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT customers_account_status_check CHECK (((account_status)::text = ANY ((ARRAY['active'::character varying, 'suspended'::character varying, 'terminated'::character varying, 'pending'::character varying])::text[])))
);


ALTER TABLE public.customers OWNER TO e173_user;

--
-- Name: customers_id_seq; Type: SEQUENCE; Schema: public; Owner: e173_user
--

CREATE SEQUENCE public.customers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.customers_id_seq OWNER TO e173_user;

--
-- Name: customers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: e173_user
--

ALTER SEQUENCE public.customers_id_seq OWNED BY public.customers.id;


--
-- Name: gateway_health_logs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.gateway_health_logs (
    id integer NOT NULL,
    gateway_id uuid,
    status character varying(50) NOT NULL,
    response_time_ms integer,
    concurrent_calls integer,
    cpu_usage numeric(5,2),
    memory_usage numeric(5,2),
    error_rate numeric(5,2),
    checked_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.gateway_health_logs OWNER TO postgres;

--
-- Name: gateway_health_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.gateway_health_logs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.gateway_health_logs_id_seq OWNER TO postgres;

--
-- Name: gateway_health_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.gateway_health_logs_id_seq OWNED BY public.gateway_health_logs.id;


--
-- Name: gateways; Type: TABLE; Schema: public; Owner: e173_user
--

CREATE TABLE public.gateways (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    location character varying(255),
    ami_host character varying(255) NOT NULL,
    ami_port character varying(10) DEFAULT '5038'::character varying NOT NULL,
    ami_user character varying(255) NOT NULL,
    ami_pass character varying(255) NOT NULL,
    status character varying(50) DEFAULT 'offline'::character varying NOT NULL,
    enabled boolean DEFAULT true NOT NULL,
    last_seen timestamp with time zone,
    last_error text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    sip_endpoint character varying(255),
    sip_port integer DEFAULT 5060,
    health_status character varying(50) DEFAULT 'healthy'::character varying,
    health_last_check timestamp without time zone DEFAULT now(),
    concurrent_calls integer DEFAULT 0,
    max_concurrent_calls integer DEFAULT 50,
    region character varying(100),
    operator_preferences jsonb
);


ALTER TABLE public.gateways OWNER TO e173_user;

--
-- Name: modems; Type: TABLE; Schema: public; Owner: e173_user
--

CREATE TABLE public.modems (
    id integer NOT NULL,
    device_path character varying(255) NOT NULL,
    imei character varying(20),
    imsi character varying(20),
    model character varying(100),
    manufacturer character varying(100),
    firmware_version character varying(100),
    signal_strength_dbm integer,
    network_operator_name character varying(100),
    network_registration_status character varying(50),
    status character varying(50) DEFAULT 'unknown'::character varying NOT NULL,
    last_seen_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.modems OWNER TO e173_user;

--
-- Name: modems_id_seq; Type: SEQUENCE; Schema: public; Owner: e173_user
--

CREATE SEQUENCE public.modems_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.modems_id_seq OWNER TO e173_user;

--
-- Name: modems_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: e173_user
--

ALTER SEQUENCE public.modems_id_seq OWNED BY public.modems.id;


--
-- Name: notification_templates; Type: TABLE; Schema: public; Owner: e173_user
--

CREATE TABLE public.notification_templates (
    id bigint NOT NULL,
    template_name character varying(100) NOT NULL,
    template_type character varying(50) NOT NULL,
    subject_template text,
    body_template text NOT NULL,
    variables jsonb,
    is_active boolean DEFAULT true,
    created_by bigint,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.notification_templates OWNER TO e173_user;

--
-- Name: notification_templates_id_seq; Type: SEQUENCE; Schema: public; Owner: e173_user
--

CREATE SEQUENCE public.notification_templates_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.notification_templates_id_seq OWNER TO e173_user;

--
-- Name: notification_templates_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: e173_user
--

ALTER SEQUENCE public.notification_templates_id_seq OWNED BY public.notification_templates.id;


--
-- Name: operator_routing_rules; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.operator_routing_rules (
    id integer NOT NULL,
    operator_name character varying(100) NOT NULL,
    prefix_pattern character varying(20) NOT NULL,
    preferred_gateway_id uuid,
    backup_gateway_ids uuid[],
    cost_per_minute numeric(8,4),
    quality_score numeric(3,2),
    active boolean DEFAULT true,
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.operator_routing_rules OWNER TO postgres;

--
-- Name: operator_routing_rules_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.operator_routing_rules_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.operator_routing_rules_id_seq OWNER TO postgres;

--
-- Name: operator_routing_rules_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.operator_routing_rules_id_seq OWNED BY public.operator_routing_rules.id;


--
-- Name: payments; Type: TABLE; Schema: public; Owner: e173_user
--

CREATE TABLE public.payments (
    id bigint NOT NULL,
    customer_id bigint NOT NULL,
    payment_reference character varying(100) NOT NULL,
    payment_type character varying(20) NOT NULL,
    amount numeric(12,4) NOT NULL,
    currency character varying(3) DEFAULT 'USD'::character varying,
    description text,
    payment_method character varying(50),
    transaction_id character varying(255),
    gateway_response jsonb,
    status character varying(20) DEFAULT 'pending'::character varying,
    processed_at timestamp with time zone,
    processed_by bigint,
    notes text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT payments_payment_type_check CHECK (((payment_type)::text = ANY ((ARRAY['credit'::character varying, 'debit'::character varying, 'refund'::character varying, 'adjustment'::character varying, 'auto_recharge'::character varying])::text[]))),
    CONSTRAINT payments_status_check CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'completed'::character varying, 'failed'::character varying, 'refunded'::character varying, 'cancelled'::character varying])::text[])))
);


ALTER TABLE public.payments OWNER TO e173_user;

--
-- Name: payments_id_seq; Type: SEQUENCE; Schema: public; Owner: e173_user
--

CREATE SEQUENCE public.payments_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.payments_id_seq OWNER TO e173_user;

--
-- Name: payments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: e173_user
--

ALTER SEQUENCE public.payments_id_seq OWNED BY public.payments.id;


--
-- Name: rate_plans; Type: TABLE; Schema: public; Owner: e173_user
--

CREATE TABLE public.rate_plans (
    id bigint NOT NULL,
    plan_name character varying(100) NOT NULL,
    plan_code character varying(20) NOT NULL,
    description text,
    currency character varying(3) DEFAULT 'USD'::character varying,
    rate_per_minute numeric(10,6) NOT NULL,
    rate_per_second numeric(10,6) NOT NULL,
    minimum_billing_seconds integer DEFAULT 1,
    connection_fee numeric(10,4) DEFAULT 0.0,
    daily_cap numeric(10,4),
    monthly_cap numeric(10,4),
    effective_from timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    effective_until timestamp with time zone,
    is_active boolean DEFAULT true,
    created_by bigint,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.rate_plans OWNER TO e173_user;

--
-- Name: rate_plans_id_seq; Type: SEQUENCE; Schema: public; Owner: e173_user
--

CREATE SEQUENCE public.rate_plans_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.rate_plans_id_seq OWNER TO e173_user;

--
-- Name: rate_plans_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: e173_user
--

ALTER SEQUENCE public.rate_plans_id_seq OWNED BY public.rate_plans.id;


--
-- Name: revenue_tracking; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.revenue_tracking (
    id integer NOT NULL,
    call_id character varying(255),
    revenue_source character varying(50),
    amount numeric(10,4) NOT NULL,
    currency character varying(3) DEFAULT 'USD'::character varying,
    customer_id integer,
    gateway_id uuid,
    ai_agent_id integer,
    billing_date date DEFAULT CURRENT_DATE,
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.revenue_tracking OWNER TO postgres;

--
-- Name: revenue_tracking_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.revenue_tracking_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.revenue_tracking_id_seq OWNER TO postgres;

--
-- Name: revenue_tracking_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.revenue_tracking_id_seq OWNED BY public.revenue_tracking.id;


--
-- Name: routing_rules; Type: TABLE; Schema: public; Owner: e173_user
--

CREATE TABLE public.routing_rules (
    id bigint NOT NULL,
    rule_name character varying(100) NOT NULL,
    rule_order integer DEFAULT 1000 NOT NULL,
    prefix_pattern character varying(50) NOT NULL,
    destination_pattern character varying(50),
    caller_id_pattern character varying(50),
    route_to_modem_id bigint,
    route_to_pool character varying(50),
    max_channels integer DEFAULT 1,
    time_restrictions jsonb,
    customer_restrictions bigint[] DEFAULT '{}'::bigint[],
    cost_markup_percent numeric(5,2) DEFAULT 0.0,
    is_active boolean DEFAULT true,
    notes text,
    created_by bigint,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.routing_rules OWNER TO e173_user;

--
-- Name: routing_rules_id_seq; Type: SEQUENCE; Schema: public; Owner: e173_user
--

CREATE SEQUENCE public.routing_rules_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.routing_rules_id_seq OWNER TO e173_user;

--
-- Name: routing_rules_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: e173_user
--

ALTER SEQUENCE public.routing_rules_id_seq OWNED BY public.routing_rules.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: e173_user
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO e173_user;

--
-- Name: sim_cards; Type: TABLE; Schema: public; Owner: e173_user
--

CREATE TABLE public.sim_cards (
    id integer NOT NULL,
    modem_id integer,
    iccid character varying(25) NOT NULL,
    imsi character varying(20),
    msisdn character varying(20),
    operator_name character varying(100),
    network_country_code character varying(10),
    balance numeric(10,4) DEFAULT 0.00,
    balance_currency character varying(10),
    balance_last_checked_at timestamp with time zone,
    data_allowance_mb integer,
    data_used_mb integer,
    status character varying(50) DEFAULT 'unknown'::character varying NOT NULL,
    pin1 character varying(10),
    puk1 character varying(10),
    pin2 character varying(10),
    puk2 character varying(10),
    activation_date date,
    expiry_date date,
    recharge_history jsonb,
    notes text,
    cell_id character varying(50),
    lac character varying(50),
    psc character varying(50),
    rscp integer,
    ecio integer,
    bts_info_history jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.sim_cards OWNER TO e173_user;

--
-- Name: sim_cards_id_seq; Type: SEQUENCE; Schema: public; Owner: e173_user
--

CREATE SEQUENCE public.sim_cards_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.sim_cards_id_seq OWNER TO e173_user;

--
-- Name: sim_cards_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: e173_user
--

ALTER SEQUENCE public.sim_cards_id_seq OWNED BY public.sim_cards.id;


--
-- Name: sim_pool_assignments; Type: TABLE; Schema: public; Owner: e173_user
--

CREATE TABLE public.sim_pool_assignments (
    id bigint NOT NULL,
    sim_pool_id bigint NOT NULL,
    sim_card_id bigint NOT NULL,
    priority integer DEFAULT 100,
    is_active boolean DEFAULT true,
    assigned_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    assigned_by bigint
);


ALTER TABLE public.sim_pool_assignments OWNER TO e173_user;

--
-- Name: sim_pool_assignments_id_seq; Type: SEQUENCE; Schema: public; Owner: e173_user
--

CREATE SEQUENCE public.sim_pool_assignments_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.sim_pool_assignments_id_seq OWNER TO e173_user;

--
-- Name: sim_pool_assignments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: e173_user
--

ALTER SEQUENCE public.sim_pool_assignments_id_seq OWNED BY public.sim_pool_assignments.id;


--
-- Name: sim_pools; Type: TABLE; Schema: public; Owner: e173_user
--

CREATE TABLE public.sim_pools (
    id bigint NOT NULL,
    pool_name character varying(50) NOT NULL,
    description text,
    load_balance_method character varying(20) DEFAULT 'round_robin'::character varying,
    max_channels_per_sim integer DEFAULT 1,
    is_active boolean DEFAULT true,
    created_by bigint,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT sim_pools_load_balance_method_check CHECK (((load_balance_method)::text = ANY ((ARRAY['round_robin'::character varying, 'least_used'::character varying, 'random'::character varying, 'failover'::character varying])::text[])))
);


ALTER TABLE public.sim_pools OWNER TO e173_user;

--
-- Name: sim_pools_id_seq; Type: SEQUENCE; Schema: public; Owner: e173_user
--

CREATE SEQUENCE public.sim_pools_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.sim_pools_id_seq OWNER TO e173_user;

--
-- Name: sim_pools_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: e173_user
--

ALTER SEQUENCE public.sim_pools_id_seq OWNED BY public.sim_pools.id;


--
-- Name: sip_calls; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.sip_calls (
    id integer NOT NULL,
    call_id character varying(255) NOT NULL,
    caller_number character varying(50) NOT NULL,
    destination_number character varying(50) NOT NULL,
    gateway_id uuid,
    filter_result jsonb,
    routed_to_ai boolean DEFAULT false,
    ai_session_id character varying(255),
    billing_seconds integer DEFAULT 0,
    spam_score numeric(3,2),
    operator_detected character varying(50),
    sticky_routing boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT now(),
    ended_at timestamp without time zone
);


ALTER TABLE public.sip_calls OWNER TO postgres;

--
-- Name: sip_calls_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.sip_calls_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.sip_calls_id_seq OWNER TO postgres;

--
-- Name: sip_calls_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.sip_calls_id_seq OWNED BY public.sip_calls.id;


--
-- Name: system_config; Type: TABLE; Schema: public; Owner: e173_user
--

CREATE TABLE public.system_config (
    id bigint NOT NULL,
    config_key character varying(100) NOT NULL,
    config_value text,
    config_type character varying(20) DEFAULT 'string'::character varying,
    description text,
    is_encrypted boolean DEFAULT false,
    is_system boolean DEFAULT false,
    category character varying(50),
    updated_by bigint,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT system_config_config_type_check CHECK (((config_type)::text = ANY ((ARRAY['string'::character varying, 'integer'::character varying, 'float'::character varying, 'boolean'::character varying, 'json'::character varying])::text[])))
);


ALTER TABLE public.system_config OWNER TO e173_user;

--
-- Name: system_config_id_seq; Type: SEQUENCE; Schema: public; Owner: e173_user
--

CREATE SEQUENCE public.system_config_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.system_config_id_seq OWNER TO e173_user;

--
-- Name: system_config_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: e173_user
--

ALTER SEQUENCE public.system_config_id_seq OWNED BY public.system_config.id;


--
-- Name: user_notifications; Type: TABLE; Schema: public; Owner: e173_user
--

CREATE TABLE public.user_notifications (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    notification_type character varying(50) NOT NULL,
    title character varying(255) NOT NULL,
    message text NOT NULL,
    priority character varying(20) DEFAULT 'normal'::character varying,
    is_read boolean DEFAULT false,
    read_at timestamp with time zone,
    action_url character varying(500),
    metadata jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT user_notifications_priority_check CHECK (((priority)::text = ANY ((ARRAY['low'::character varying, 'normal'::character varying, 'high'::character varying, 'critical'::character varying])::text[])))
);


ALTER TABLE public.user_notifications OWNER TO e173_user;

--
-- Name: user_notifications_id_seq; Type: SEQUENCE; Schema: public; Owner: e173_user
--

CREATE SEQUENCE public.user_notifications_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_notifications_id_seq OWNER TO e173_user;

--
-- Name: user_notifications_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: e173_user
--

ALTER SEQUENCE public.user_notifications_id_seq OWNED BY public.user_notifications.id;


--
-- Name: user_sessions; Type: TABLE; Schema: public; Owner: e173_user
--

CREATE TABLE public.user_sessions (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    session_token character varying(255) NOT NULL,
    ip_address inet,
    user_agent text,
    expires_at timestamp with time zone NOT NULL,
    is_active boolean DEFAULT true,
    last_activity_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.user_sessions OWNER TO e173_user;

--
-- Name: user_sessions_id_seq; Type: SEQUENCE; Schema: public; Owner: e173_user
--

CREATE SEQUENCE public.user_sessions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_sessions_id_seq OWNER TO e173_user;

--
-- Name: user_sessions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: e173_user
--

ALTER SEQUENCE public.user_sessions_id_seq OWNED BY public.user_sessions.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: e173_user
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    username character varying(50) NOT NULL,
    email character varying(255) NOT NULL,
    password_hash character varying(255) NOT NULL,
    first_name character varying(100),
    last_name character varying(100),
    role character varying(20) NOT NULL,
    is_active boolean DEFAULT true,
    is_2fa_enabled boolean DEFAULT false,
    two_fa_secret character varying(32),
    last_login_at timestamp with time zone,
    failed_login_attempts integer DEFAULT 0,
    locked_until timestamp with time zone,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_role_check CHECK (((role)::text = ANY ((ARRAY['super_admin'::character varying, 'admin'::character varying, 'manager'::character varying, 'employee'::character varying, 'viewer'::character varying])::text[])))
);


ALTER TABLE public.users OWNER TO e173_user;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: e173_user
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO e173_user;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: e173_user
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: whatsapp_validation_cache; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.whatsapp_validation_cache (
    phone_number character varying(50) NOT NULL,
    has_whatsapp boolean NOT NULL,
    confidence_score numeric(3,2),
    checked_at timestamp without time zone DEFAULT now(),
    expires_at timestamp without time zone DEFAULT (now() + '7 days'::interval)
);


ALTER TABLE public.whatsapp_validation_cache OWNER TO postgres;

--
-- Name: ai_voice_agents id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_voice_agents ALTER COLUMN id SET DEFAULT nextval('public.ai_voice_agents_id_seq'::regclass);


--
-- Name: ai_voice_sessions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_voice_sessions ALTER COLUMN id SET DEFAULT nextval('public.ai_voice_sessions_id_seq'::regclass);


--
-- Name: audit_logs id; Type: DEFAULT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.audit_logs ALTER COLUMN id SET DEFAULT nextval('public.audit_logs_id_seq'::regclass);


--
-- Name: blacklist id; Type: DEFAULT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.blacklist ALTER COLUMN id SET DEFAULT nextval('public.blacklist_id_seq'::regclass);


--
-- Name: call_detail_records id; Type: DEFAULT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.call_detail_records ALTER COLUMN id SET DEFAULT nextval('public.call_detail_records_id_seq'::regclass);


--
-- Name: call_patterns id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.call_patterns ALTER COLUMN id SET DEFAULT nextval('public.call_patterns_id_seq'::regclass);


--
-- Name: customer_rate_plans id; Type: DEFAULT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.customer_rate_plans ALTER COLUMN id SET DEFAULT nextval('public.customer_rate_plans_id_seq'::regclass);


--
-- Name: customers id; Type: DEFAULT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.customers ALTER COLUMN id SET DEFAULT nextval('public.customers_id_seq'::regclass);


--
-- Name: gateway_health_logs id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.gateway_health_logs ALTER COLUMN id SET DEFAULT nextval('public.gateway_health_logs_id_seq'::regclass);


--
-- Name: modems id; Type: DEFAULT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.modems ALTER COLUMN id SET DEFAULT nextval('public.modems_id_seq'::regclass);


--
-- Name: notification_templates id; Type: DEFAULT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.notification_templates ALTER COLUMN id SET DEFAULT nextval('public.notification_templates_id_seq'::regclass);


--
-- Name: operator_routing_rules id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.operator_routing_rules ALTER COLUMN id SET DEFAULT nextval('public.operator_routing_rules_id_seq'::regclass);


--
-- Name: payments id; Type: DEFAULT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.payments ALTER COLUMN id SET DEFAULT nextval('public.payments_id_seq'::regclass);


--
-- Name: rate_plans id; Type: DEFAULT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.rate_plans ALTER COLUMN id SET DEFAULT nextval('public.rate_plans_id_seq'::regclass);


--
-- Name: revenue_tracking id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.revenue_tracking ALTER COLUMN id SET DEFAULT nextval('public.revenue_tracking_id_seq'::regclass);


--
-- Name: routing_rules id; Type: DEFAULT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.routing_rules ALTER COLUMN id SET DEFAULT nextval('public.routing_rules_id_seq'::regclass);


--
-- Name: sim_cards id; Type: DEFAULT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.sim_cards ALTER COLUMN id SET DEFAULT nextval('public.sim_cards_id_seq'::regclass);


--
-- Name: sim_pool_assignments id; Type: DEFAULT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.sim_pool_assignments ALTER COLUMN id SET DEFAULT nextval('public.sim_pool_assignments_id_seq'::regclass);


--
-- Name: sim_pools id; Type: DEFAULT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.sim_pools ALTER COLUMN id SET DEFAULT nextval('public.sim_pools_id_seq'::regclass);


--
-- Name: sip_calls id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sip_calls ALTER COLUMN id SET DEFAULT nextval('public.sip_calls_id_seq'::regclass);


--
-- Name: system_config id; Type: DEFAULT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.system_config ALTER COLUMN id SET DEFAULT nextval('public.system_config_id_seq'::regclass);


--
-- Name: user_notifications id; Type: DEFAULT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.user_notifications ALTER COLUMN id SET DEFAULT nextval('public.user_notifications_id_seq'::regclass);


--
-- Name: user_sessions id; Type: DEFAULT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.user_sessions ALTER COLUMN id SET DEFAULT nextval('public.user_sessions_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: ai_voice_agents; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.ai_voice_agents (id, agent_name, endpoint_url, api_key, model_type, voice_id, language, concurrent_limit, current_sessions, cost_per_minute, quality_rating, active, created_at) FROM stdin;
1	SpamBot Interceptor 1	http://ai-service:8080/voice/session1	\N	elevenlabs	voice_001	en	10	0	0.0200	0.90	t	2025-06-28 01:02:55.109994
2	SpamBot Interceptor 2	http://ai-service:8080/voice/session2	\N	elevenlabs	voice_002	en	10	0	0.0200	0.88	t	2025-06-28 01:02:55.109994
3	Customer Service AI	http://ai-service:8080/voice/customer	\N	openai	voice_customer	en	10	0	0.0500	0.95	t	2025-06-28 01:02:55.109994
4	SpamBot Interceptor 1	http://ai-service:8080/voice/session1	\N	elevenlabs	voice_001	en	10	0	0.0200	0.90	t	2025-06-28 01:06:38.678684
5	SpamBot Interceptor 2	http://ai-service:8080/voice/session2	\N	elevenlabs	voice_002	en	10	0	0.0200	0.88	t	2025-06-28 01:06:38.678684
6	Customer Service AI	http://ai-service:8080/voice/customer	\N	openai	voice_customer	en	10	0	0.0500	0.95	t	2025-06-28 01:06:38.678684
\.


--
-- Data for Name: ai_voice_sessions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.ai_voice_sessions (id, session_id, agent_id, sip_call_id, caller_number, session_duration, tokens_used, revenue_generated, conversation_log, started_at, ended_at) FROM stdin;
\.


--
-- Data for Name: audit_logs; Type: TABLE DATA; Schema: public; Owner: e173_user
--

COPY public.audit_logs (id, user_id, action, entity_type, entity_id, old_values, new_values, ip_address, user_agent, success, error_message, created_at) FROM stdin;
1	1	login	\N	\N	\N	\N	127.0.0.1	curl/7.81.0	t	\N	2025-06-25 00:57:06.133691+02
2	1	login	\N	\N	\N	\N	192.168.1.41	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Safari/605.1.15	t	\N	2025-06-25 00:57:53.138057+02
3	1	login	\N	\N	\N	\N	192.168.1.41	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Safari/605.1.15	t	\N	2025-06-25 00:58:24.749918+02
4	1	login	\N	\N	\N	\N	127.0.0.1	curl/7.81.0	t	\N	2025-06-25 00:58:54.103906+02
5	1	login	\N	\N	\N	\N	192.168.1.41	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Safari/605.1.15	t	\N	2025-06-25 01:00:35.217463+02
6	1	login	\N	\N	\N	\N	192.168.1.41	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Safari/605.1.15	t	\N	2025-06-25 01:01:41.908534+02
7	1	login	\N	\N	\N	\N	127.0.0.1	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36	t	\N	2025-06-25 01:04:16.860065+02
8	1	login	\N	\N	\N	\N	127.0.0.1	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Code/1.101.1 Chrome/134.0.6998.205 Electron/35.5.1 Safari/537.36	t	\N	2025-06-26 02:44:51.378118+02
41	1	login	\N	\N	\N	\N	192.168.1.41	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Safari/605.1.15	t	\N	2025-06-28 01:18:34.12046+02
42	1	login	\N	\N	\N	\N	127.0.0.1	curl/7.81.0	t	\N	2025-06-28 01:31:44.170848+02
\.


--
-- Data for Name: blacklist; Type: TABLE DATA; Schema: public; Owner: e173_user
--

COPY public.blacklist (id, number_pattern, blacklist_type, reason, auto_added, detection_method, block_inbound, block_outbound, temporary_until, violation_count, last_violation_at, created_by, created_at, updated_at, reason_code, expires_at, confidence_score) FROM stdin;
\.


--
-- Data for Name: call_detail_records; Type: TABLE DATA; Schema: public; Owner: e173_user
--

COPY public.call_detail_records (id, sim_card_id, modem_id, asterisk_unique_id, call_direction, source_number, destination_number, call_start_time, call_answer_time, call_end_time, duration_seconds, billable_duration_seconds, disposition, hangup_cause, cost_per_minute, total_cost, customer_id, is_spam, spam_reason, recorded_audio_path, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: call_patterns; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.call_patterns (id, phone_number, total_calls, answered_calls, short_calls, avg_call_duration, spam_score, last_call_time, pattern_updated) FROM stdin;
\.


--
-- Data for Name: customer_rate_plans; Type: TABLE DATA; Schema: public; Owner: e173_user
--

COPY public.customer_rate_plans (id, customer_id, rate_plan_id, effective_from, effective_until, is_active, created_by, created_at) FROM stdin;
\.


--
-- Data for Name: customers; Type: TABLE DATA; Schema: public; Owner: e173_user
--

COPY public.customers (id, customer_code, company_name, contact_person, email, phone, address, city, state, country, postal_code, billing_address, billing_city, billing_state, billing_country, billing_postal_code, account_status, credit_limit, current_balance, monthly_limit, timezone, preferred_currency, auto_recharge_enabled, auto_recharge_threshold, auto_recharge_amount, notes, created_by, assigned_to, created_at, updated_at) FROM stdin;
4	CUST004	Suspended Corp	Jane Doe	jane@suspended.com	+1-555-0103	321 Inactive St	Chicago	IL	USA	60601	\N	\N	\N	\N	\N	suspended	1000.0000	-25.0000	500.0000	America/Chicago	USD	f	\N	\N	Suspended due to negative balance	\N	\N	2025-06-24 17:41:20.519629+02	2025-06-24 17:41:20.519629+02
1	CUST001	TechCorp Solutions	John Smith	john@techcorp.com	+1-555-0101	123 Business Ave	New York	NY	USA	10001	\N	\N	\N	\N	\N	active	1000.0000	350.7500	500.0000	America/New_York	USD	t	50.0000	100.0000	Premium customer with auto-recharge	\N	\N	2025-06-24 17:41:20.519629+02	2025-06-24 17:41:20.526104+02
2	CUST002	Global Communications Ltd	Sarah Johnson	sarah@globalcomm.com	+44-20-7946-0958	456 International Blvd	London	England	UK	SW1A 1AA	\N	\N	\N	\N	\N	active	2000.0000	350.2500	1000.0000	Europe/London	GBP	f	\N	\N	Enterprise customer - manual billing	\N	\N	2025-06-24 17:41:20.519629+02	2025-06-24 17:41:20.531077+02
3	CUST003	StartupTech Inc	Mike Chen	mike@startuptech.com	+1-555-0102	789 Innovation Dr	San Francisco	CA	USA	94105	\N	\N	\N	\N	\N	active	500.0000	125.5000	250.0000	America/Los_Angeles	USD	t	25.0000	50.0000	Growing startup customer	\N	\N	2025-06-24 17:41:20.519629+02	2025-06-24 17:41:20.531754+02
\.


--
-- Data for Name: gateway_health_logs; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.gateway_health_logs (id, gateway_id, status, response_time_ms, concurrent_calls, cpu_usage, memory_usage, error_rate, checked_at) FROM stdin;
\.


--
-- Data for Name: gateways; Type: TABLE DATA; Schema: public; Owner: e173_user
--

COPY public.gateways (id, name, description, location, ami_host, ami_port, ami_user, ami_pass, status, enabled, last_seen, last_error, created_at, updated_at, sip_endpoint, sip_port, health_status, health_last_check, concurrent_calls, max_concurrent_calls, region, operator_preferences) FROM stdin;
\.


--
-- Data for Name: modems; Type: TABLE DATA; Schema: public; Owner: e173_user
--

COPY public.modems (id, device_path, imei, imsi, model, manufacturer, firmware_version, signal_strength_dbm, network_operator_name, network_registration_status, status, last_seen_at, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: notification_templates; Type: TABLE DATA; Schema: public; Owner: e173_user
--

COPY public.notification_templates (id, template_name, template_type, subject_template, body_template, variables, is_active, created_by, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: operator_routing_rules; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.operator_routing_rules (id, operator_name, prefix_pattern, preferred_gateway_id, backup_gateway_ids, cost_per_minute, quality_score, active, created_at) FROM stdin;
\.


--
-- Data for Name: payments; Type: TABLE DATA; Schema: public; Owner: e173_user
--

COPY public.payments (id, customer_id, payment_reference, payment_type, amount, currency, description, payment_method, transaction_id, gateway_response, status, processed_at, processed_by, notes, created_at, updated_at) FROM stdin;
2	1	CC-2024-001	credit	100.0000	USD	\N	credit_card	\N	\N	completed	\N	1	Monthly auto-recharge	2025-06-24 17:41:58.483314+02	2025-06-24 17:41:58.483314+02
3	2	BT-2024-001	credit	200.0000	USD	\N	bank_transfer	\N	\N	completed	\N	1	Manual payment received	2025-06-24 17:41:58.483314+02	2025-06-24 17:41:58.483314+02
4	3	CC-2024-002	credit	50.0000	USD	\N	credit_card	\N	\N	completed	\N	1	Auto-recharge payment	2025-06-24 17:41:58.483314+02	2025-06-24 17:41:58.483314+02
\.


--
-- Data for Name: rate_plans; Type: TABLE DATA; Schema: public; Owner: e173_user
--

COPY public.rate_plans (id, plan_name, plan_code, description, currency, rate_per_minute, rate_per_second, minimum_billing_seconds, connection_fee, daily_cap, monthly_cap, effective_from, effective_until, is_active, created_by, created_at, updated_at) FROM stdin;
1	Standard Plan	STD001	Default standard rate plan	USD	0.050000	0.000833	1	0.0000	\N	\N	2025-06-22 21:53:03.513731+02	\N	t	\N	2025-06-22 21:53:03.513731+02	2025-06-22 21:53:03.513731+02
\.


--
-- Data for Name: revenue_tracking; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.revenue_tracking (id, call_id, revenue_source, amount, currency, customer_id, gateway_id, ai_agent_id, billing_date, created_at) FROM stdin;
\.


--
-- Data for Name: routing_rules; Type: TABLE DATA; Schema: public; Owner: e173_user
--

COPY public.routing_rules (id, rule_name, rule_order, prefix_pattern, destination_pattern, caller_id_pattern, route_to_modem_id, route_to_pool, max_channels, time_restrictions, customer_restrictions, cost_markup_percent, is_active, notes, created_by, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: e173_user
--

COPY public.schema_migrations (version, dirty) FROM stdin;
10	f
\.


--
-- Data for Name: sim_cards; Type: TABLE DATA; Schema: public; Owner: e173_user
--

COPY public.sim_cards (id, modem_id, iccid, imsi, msisdn, operator_name, network_country_code, balance, balance_currency, balance_last_checked_at, data_allowance_mb, data_used_mb, status, pin1, puk1, pin2, puk2, activation_date, expiry_date, recharge_history, notes, cell_id, lac, psc, rscp, ecio, bts_info_history, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: sim_pool_assignments; Type: TABLE DATA; Schema: public; Owner: e173_user
--

COPY public.sim_pool_assignments (id, sim_pool_id, sim_card_id, priority, is_active, assigned_at, assigned_by) FROM stdin;
\.


--
-- Data for Name: sim_pools; Type: TABLE DATA; Schema: public; Owner: e173_user
--

COPY public.sim_pools (id, pool_name, description, load_balance_method, max_channels_per_sim, is_active, created_by, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: sip_calls; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.sip_calls (id, call_id, caller_number, destination_number, gateway_id, filter_result, routed_to_ai, ai_session_id, billing_seconds, spam_score, operator_detected, sticky_routing, created_at, ended_at) FROM stdin;
\.


--
-- Data for Name: system_config; Type: TABLE DATA; Schema: public; Owner: e173_user
--

COPY public.system_config (id, config_key, config_value, config_type, description, is_encrypted, is_system, category, updated_by, created_at, updated_at) FROM stdin;
1	system_name	E173 Gateway Management System	string	System display name	f	f	general	\N	2025-06-22 21:53:03.677402+02	2025-06-22 21:53:03.677402+02
2	default_timezone	UTC	string	Default system timezone	f	f	general	\N	2025-06-22 21:53:03.677402+02	2025-06-22 21:53:03.677402+02
3	session_timeout_minutes	1440	integer	User session timeout in minutes	f	f	security	\N	2025-06-22 21:53:03.677402+02	2025-06-22 21:53:03.677402+02
4	max_failed_login_attempts	5	integer	Maximum failed login attempts before account lock	f	f	security	\N	2025-06-22 21:53:03.677402+02	2025-06-22 21:53:03.677402+02
5	account_lockout_minutes	30	integer	Account lockout duration in minutes	f	f	security	\N	2025-06-22 21:53:03.677402+02	2025-06-22 21:53:03.677402+02
6	spam_detection_enabled	true	boolean	Enable automatic spam detection	f	f	routing	\N	2025-06-22 21:53:03.677402+02	2025-06-22 21:53:03.677402+02
7	spam_short_call_threshold	10	integer	Calls shorter than this (seconds) are flagged as potential spam	f	f	routing	\N	2025-06-22 21:53:03.677402+02	2025-06-22 21:53:03.677402+02
8	auto_blacklist_enabled	true	boolean	Automatically blacklist numbers flagged as spam	f	f	routing	\N	2025-06-22 21:53:03.677402+02	2025-06-22 21:53:03.677402+02
9	low_balance_threshold	10.00	float	Send alerts when customer balance falls below this amount	f	f	billing	\N	2025-06-22 21:53:03.677402+02	2025-06-22 21:53:03.677402+02
10	currency_symbol	$	string	Default currency symbol	f	f	billing	\N	2025-06-22 21:53:03.677402+02	2025-06-22 21:53:03.677402+02
\.


--
-- Data for Name: user_notifications; Type: TABLE DATA; Schema: public; Owner: e173_user
--

COPY public.user_notifications (id, user_id, notification_type, title, message, priority, is_read, read_at, action_url, metadata, created_at) FROM stdin;
\.


--
-- Data for Name: user_sessions; Type: TABLE DATA; Schema: public; Owner: e173_user
--

COPY public.user_sessions (id, user_id, session_token, ip_address, user_agent, expires_at, is_active, last_activity_at, created_at) FROM stdin;
1	1	f94121cb15aaee9425c3758f7b21e4d87bb7ea292785b058cf9d2f04b345db4a	127.0.0.1	curl/7.81.0	2025-06-26 00:57:06.128056+02	t	2025-06-25 00:57:25.094794+02	2025-06-25 00:57:06.128986+02
2	1	2fbbfa35bcd547160e9d6aefeb87e9e3030835ebcffe40d9ced42a3f088535da	192.168.1.41	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Safari/605.1.15	2025-06-26 00:57:53.136155+02	t	2025-06-25 00:57:53.275476+02	2025-06-25 00:57:53.136349+02
4	1	ad145f3343aea14ab7df4234887fd614b3c705e3c892e55fd8f7f5fcb660e839	127.0.0.1	curl/7.81.0	2025-06-26 00:58:54.10187+02	t	2025-06-25 00:58:54.102043+02	2025-06-25 00:58:54.102043+02
3	1	7a8d5269d3bb40600e92faf2cca0b6b2318d7b8ad2122f1c9d57044fb0a97966	192.168.1.41	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Safari/605.1.15	2025-06-26 00:58:24.748033+02	t	2025-06-25 01:00:07.228187+02	2025-06-25 00:58:24.748118+02
5	1	040bc3335af86e6b4432358e702a2a4b65477d8e0dc5c250c31610f3edb16fb2	192.168.1.41	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Safari/605.1.15	2025-06-26 01:00:35.215125+02	t	2025-06-25 01:01:30.500344+02	2025-06-25 01:00:35.215281+02
6	1	8ae4f33f393bbe9143e74aabe5ce3a7b1f9de14ebdd3e12d2c37579aeab9b915	192.168.1.41	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Safari/605.1.15	2025-06-26 01:01:41.905776+02	t	2025-06-25 01:01:48.296714+02	2025-06-25 01:01:41.905912+02
7	1	f7f1eb3d76e2074a64b22dcb7fbec24daabd0fd356298a60b28f1fc981820e7c	127.0.0.1	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36	2025-06-26 01:04:16.858146+02	t	2025-06-25 01:04:16.868975+02	2025-06-25 01:04:16.858319+02
8	1	2922acde7c141891ea31a9842ccd10db661423b1e83d5b7afb509318d5adb3ca	127.0.0.1	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Code/1.101.1 Chrome/134.0.6998.205 Electron/35.5.1 Safari/537.36	2025-06-27 02:44:51.359495+02	t	2025-06-26 02:44:51.360914+02	2025-06-26 02:44:51.360914+02
41	1	7b15933bccd5fd05619d23f9469d80e49b092f47b3b01b17ceb83009409c5430	192.168.1.41	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Safari/605.1.15	2025-06-29 01:18:34.113311+02	t	2025-06-28 01:29:50.17348+02	2025-06-28 01:18:34.113884+02
42	1	b144860f219aacc238423e040814b900cede8fb89f62d1af7abe389b9665ab73	127.0.0.1	curl/7.81.0	2025-06-29 01:31:44.168787+02	t	2025-06-28 01:31:48.57988+02	2025-06-28 01:31:44.16894+02
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: e173_user
--

COPY public.users (id, username, email, password_hash, first_name, last_name, role, is_active, is_2fa_enabled, two_fa_secret, last_login_at, failed_login_attempts, locked_until, created_at, updated_at) FROM stdin;
1	admin	admin@e173gateway.local	$2b$12$qRiY5cNnl44T8Tf0l1W.7uYJ5cYIliQ555Con4k8av7LH.Y76ktcm	System	Administrator	super_admin	t	f	\N	2025-06-28 01:31:44.170081+02	0	\N	2025-06-22 21:51:57.310784+02	2025-06-28 01:31:44.170081+02
\.


--
-- Data for Name: whatsapp_validation_cache; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.whatsapp_validation_cache (phone_number, has_whatsapp, confidence_score, checked_at, expires_at) FROM stdin;
\.


--
-- Name: ai_voice_agents_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ai_voice_agents_id_seq', 6, true);


--
-- Name: ai_voice_sessions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ai_voice_sessions_id_seq', 1, false);


--
-- Name: audit_logs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: e173_user
--

SELECT pg_catalog.setval('public.audit_logs_id_seq', 42, true);


--
-- Name: blacklist_id_seq; Type: SEQUENCE SET; Schema: public; Owner: e173_user
--

SELECT pg_catalog.setval('public.blacklist_id_seq', 1, false);


--
-- Name: call_detail_records_id_seq; Type: SEQUENCE SET; Schema: public; Owner: e173_user
--

SELECT pg_catalog.setval('public.call_detail_records_id_seq', 1, false);


--
-- Name: call_patterns_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.call_patterns_id_seq', 1, false);


--
-- Name: customer_rate_plans_id_seq; Type: SEQUENCE SET; Schema: public; Owner: e173_user
--

SELECT pg_catalog.setval('public.customer_rate_plans_id_seq', 1, false);


--
-- Name: customers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: e173_user
--

SELECT pg_catalog.setval('public.customers_id_seq', 4, true);


--
-- Name: gateway_health_logs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.gateway_health_logs_id_seq', 1, false);


--
-- Name: modems_id_seq; Type: SEQUENCE SET; Schema: public; Owner: e173_user
--

SELECT pg_catalog.setval('public.modems_id_seq', 1, false);


--
-- Name: notification_templates_id_seq; Type: SEQUENCE SET; Schema: public; Owner: e173_user
--

SELECT pg_catalog.setval('public.notification_templates_id_seq', 1, false);


--
-- Name: operator_routing_rules_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.operator_routing_rules_id_seq', 1, false);


--
-- Name: payments_id_seq; Type: SEQUENCE SET; Schema: public; Owner: e173_user
--

SELECT pg_catalog.setval('public.payments_id_seq', 4, true);


--
-- Name: rate_plans_id_seq; Type: SEQUENCE SET; Schema: public; Owner: e173_user
--

SELECT pg_catalog.setval('public.rate_plans_id_seq', 1, true);


--
-- Name: revenue_tracking_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.revenue_tracking_id_seq', 1, false);


--
-- Name: routing_rules_id_seq; Type: SEQUENCE SET; Schema: public; Owner: e173_user
--

SELECT pg_catalog.setval('public.routing_rules_id_seq', 1, false);


--
-- Name: sim_cards_id_seq; Type: SEQUENCE SET; Schema: public; Owner: e173_user
--

SELECT pg_catalog.setval('public.sim_cards_id_seq', 1, false);


--
-- Name: sim_pool_assignments_id_seq; Type: SEQUENCE SET; Schema: public; Owner: e173_user
--

SELECT pg_catalog.setval('public.sim_pool_assignments_id_seq', 1, false);


--
-- Name: sim_pools_id_seq; Type: SEQUENCE SET; Schema: public; Owner: e173_user
--

SELECT pg_catalog.setval('public.sim_pools_id_seq', 1, false);


--
-- Name: sip_calls_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.sip_calls_id_seq', 1, false);


--
-- Name: system_config_id_seq; Type: SEQUENCE SET; Schema: public; Owner: e173_user
--

SELECT pg_catalog.setval('public.system_config_id_seq', 10, true);


--
-- Name: user_notifications_id_seq; Type: SEQUENCE SET; Schema: public; Owner: e173_user
--

SELECT pg_catalog.setval('public.user_notifications_id_seq', 1, false);


--
-- Name: user_sessions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: e173_user
--

SELECT pg_catalog.setval('public.user_sessions_id_seq', 42, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: e173_user
--

SELECT pg_catalog.setval('public.users_id_seq', 1, true);


--
-- Name: ai_voice_agents ai_voice_agents_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_voice_agents
    ADD CONSTRAINT ai_voice_agents_pkey PRIMARY KEY (id);


--
-- Name: ai_voice_sessions ai_voice_sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_voice_sessions
    ADD CONSTRAINT ai_voice_sessions_pkey PRIMARY KEY (id);


--
-- Name: ai_voice_sessions ai_voice_sessions_session_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_voice_sessions
    ADD CONSTRAINT ai_voice_sessions_session_id_key UNIQUE (session_id);


--
-- Name: audit_logs audit_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.audit_logs
    ADD CONSTRAINT audit_logs_pkey PRIMARY KEY (id);


--
-- Name: blacklist blacklist_pkey; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.blacklist
    ADD CONSTRAINT blacklist_pkey PRIMARY KEY (id);


--
-- Name: call_detail_records call_detail_records_asterisk_unique_id_key; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.call_detail_records
    ADD CONSTRAINT call_detail_records_asterisk_unique_id_key UNIQUE (asterisk_unique_id);


--
-- Name: call_detail_records call_detail_records_pkey; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.call_detail_records
    ADD CONSTRAINT call_detail_records_pkey PRIMARY KEY (id);


--
-- Name: call_patterns call_patterns_phone_number_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.call_patterns
    ADD CONSTRAINT call_patterns_phone_number_key UNIQUE (phone_number);


--
-- Name: call_patterns call_patterns_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.call_patterns
    ADD CONSTRAINT call_patterns_pkey PRIMARY KEY (id);


--
-- Name: customer_rate_plans customer_rate_plans_customer_id_rate_plan_id_effective_from_key; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.customer_rate_plans
    ADD CONSTRAINT customer_rate_plans_customer_id_rate_plan_id_effective_from_key UNIQUE (customer_id, rate_plan_id, effective_from);


--
-- Name: customer_rate_plans customer_rate_plans_pkey; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.customer_rate_plans
    ADD CONSTRAINT customer_rate_plans_pkey PRIMARY KEY (id);


--
-- Name: customers customers_customer_code_key; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_customer_code_key UNIQUE (customer_code);


--
-- Name: customers customers_pkey; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_pkey PRIMARY KEY (id);


--
-- Name: gateway_health_logs gateway_health_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.gateway_health_logs
    ADD CONSTRAINT gateway_health_logs_pkey PRIMARY KEY (id);


--
-- Name: gateways gateways_pkey; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.gateways
    ADD CONSTRAINT gateways_pkey PRIMARY KEY (id);


--
-- Name: modems modems_device_path_key; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.modems
    ADD CONSTRAINT modems_device_path_key UNIQUE (device_path);


--
-- Name: modems modems_imei_key; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.modems
    ADD CONSTRAINT modems_imei_key UNIQUE (imei);


--
-- Name: modems modems_imsi_key; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.modems
    ADD CONSTRAINT modems_imsi_key UNIQUE (imsi);


--
-- Name: modems modems_pkey; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.modems
    ADD CONSTRAINT modems_pkey PRIMARY KEY (id);


--
-- Name: notification_templates notification_templates_pkey; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.notification_templates
    ADD CONSTRAINT notification_templates_pkey PRIMARY KEY (id);


--
-- Name: notification_templates notification_templates_template_name_key; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.notification_templates
    ADD CONSTRAINT notification_templates_template_name_key UNIQUE (template_name);


--
-- Name: operator_routing_rules operator_routing_rules_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.operator_routing_rules
    ADD CONSTRAINT operator_routing_rules_pkey PRIMARY KEY (id);


--
-- Name: payments payments_payment_reference_key; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_payment_reference_key UNIQUE (payment_reference);


--
-- Name: payments payments_pkey; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_pkey PRIMARY KEY (id);


--
-- Name: rate_plans rate_plans_pkey; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.rate_plans
    ADD CONSTRAINT rate_plans_pkey PRIMARY KEY (id);


--
-- Name: rate_plans rate_plans_plan_code_key; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.rate_plans
    ADD CONSTRAINT rate_plans_plan_code_key UNIQUE (plan_code);


--
-- Name: revenue_tracking revenue_tracking_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.revenue_tracking
    ADD CONSTRAINT revenue_tracking_pkey PRIMARY KEY (id);


--
-- Name: routing_rules routing_rules_pkey; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.routing_rules
    ADD CONSTRAINT routing_rules_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: sim_cards sim_cards_iccid_key; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.sim_cards
    ADD CONSTRAINT sim_cards_iccid_key UNIQUE (iccid);


--
-- Name: sim_cards sim_cards_imsi_key; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.sim_cards
    ADD CONSTRAINT sim_cards_imsi_key UNIQUE (imsi);


--
-- Name: sim_cards sim_cards_msisdn_key; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.sim_cards
    ADD CONSTRAINT sim_cards_msisdn_key UNIQUE (msisdn);


--
-- Name: sim_cards sim_cards_pkey; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.sim_cards
    ADD CONSTRAINT sim_cards_pkey PRIMARY KEY (id);


--
-- Name: sim_pool_assignments sim_pool_assignments_pkey; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.sim_pool_assignments
    ADD CONSTRAINT sim_pool_assignments_pkey PRIMARY KEY (id);


--
-- Name: sim_pool_assignments sim_pool_assignments_sim_pool_id_sim_card_id_key; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.sim_pool_assignments
    ADD CONSTRAINT sim_pool_assignments_sim_pool_id_sim_card_id_key UNIQUE (sim_pool_id, sim_card_id);


--
-- Name: sim_pools sim_pools_pkey; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.sim_pools
    ADD CONSTRAINT sim_pools_pkey PRIMARY KEY (id);


--
-- Name: sim_pools sim_pools_pool_name_key; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.sim_pools
    ADD CONSTRAINT sim_pools_pool_name_key UNIQUE (pool_name);


--
-- Name: sip_calls sip_calls_call_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sip_calls
    ADD CONSTRAINT sip_calls_call_id_key UNIQUE (call_id);


--
-- Name: sip_calls sip_calls_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sip_calls
    ADD CONSTRAINT sip_calls_pkey PRIMARY KEY (id);


--
-- Name: system_config system_config_config_key_key; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.system_config
    ADD CONSTRAINT system_config_config_key_key UNIQUE (config_key);


--
-- Name: system_config system_config_pkey; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.system_config
    ADD CONSTRAINT system_config_pkey PRIMARY KEY (id);


--
-- Name: gateways unique_gateway_name; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.gateways
    ADD CONSTRAINT unique_gateway_name UNIQUE (name);


--
-- Name: user_notifications user_notifications_pkey; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.user_notifications
    ADD CONSTRAINT user_notifications_pkey PRIMARY KEY (id);


--
-- Name: user_sessions user_sessions_pkey; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.user_sessions
    ADD CONSTRAINT user_sessions_pkey PRIMARY KEY (id);


--
-- Name: user_sessions user_sessions_session_token_key; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.user_sessions
    ADD CONSTRAINT user_sessions_session_token_key UNIQUE (session_token);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: whatsapp_validation_cache whatsapp_validation_cache_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.whatsapp_validation_cache
    ADD CONSTRAINT whatsapp_validation_cache_pkey PRIMARY KEY (phone_number);


--
-- Name: idx_ai_sessions_agent; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ai_sessions_agent ON public.ai_voice_sessions USING btree (agent_id);


--
-- Name: idx_ai_sessions_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ai_sessions_id ON public.ai_voice_sessions USING btree (session_id);


--
-- Name: idx_ai_sessions_start; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_ai_sessions_start ON public.ai_voice_sessions USING btree (started_at);


--
-- Name: idx_audit_logs_action; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_audit_logs_action ON public.audit_logs USING btree (action);


--
-- Name: idx_audit_logs_created_at; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_audit_logs_created_at ON public.audit_logs USING btree (created_at);


--
-- Name: idx_audit_logs_entity_type; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_audit_logs_entity_type ON public.audit_logs USING btree (entity_type);


--
-- Name: idx_audit_logs_user_id; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_audit_logs_user_id ON public.audit_logs USING btree (user_id);


--
-- Name: idx_blacklist_auto_added; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_blacklist_auto_added ON public.blacklist USING btree (auto_added);


--
-- Name: idx_blacklist_blacklist_type; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_blacklist_blacklist_type ON public.blacklist USING btree (blacklist_type);


--
-- Name: idx_blacklist_number_pattern; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_blacklist_number_pattern ON public.blacklist USING btree (number_pattern);


--
-- Name: idx_call_patterns_lastcall; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_call_patterns_lastcall ON public.call_patterns USING btree (last_call_time);


--
-- Name: idx_call_patterns_phone; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_call_patterns_phone ON public.call_patterns USING btree (phone_number);


--
-- Name: idx_call_patterns_spam; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_call_patterns_spam ON public.call_patterns USING btree (spam_score);


--
-- Name: idx_cdr_asterisk_unique_id; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_cdr_asterisk_unique_id ON public.call_detail_records USING btree (asterisk_unique_id);


--
-- Name: idx_cdr_call_start_time; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_cdr_call_start_time ON public.call_detail_records USING btree (call_start_time);


--
-- Name: idx_cdr_customer_id; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_cdr_customer_id ON public.call_detail_records USING btree (customer_id);


--
-- Name: idx_cdr_destination_number; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_cdr_destination_number ON public.call_detail_records USING btree (destination_number);


--
-- Name: idx_cdr_is_spam; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_cdr_is_spam ON public.call_detail_records USING btree (is_spam);


--
-- Name: idx_cdr_modem_id; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_cdr_modem_id ON public.call_detail_records USING btree (modem_id);


--
-- Name: idx_cdr_sim_card_id; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_cdr_sim_card_id ON public.call_detail_records USING btree (sim_card_id);


--
-- Name: idx_customer_rate_plans_customer_id; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_customer_rate_plans_customer_id ON public.customer_rate_plans USING btree (customer_id);


--
-- Name: idx_customer_rate_plans_is_active; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_customer_rate_plans_is_active ON public.customer_rate_plans USING btree (is_active);


--
-- Name: idx_customer_rate_plans_rate_plan_id; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_customer_rate_plans_rate_plan_id ON public.customer_rate_plans USING btree (rate_plan_id);


--
-- Name: idx_customers_account_status; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_customers_account_status ON public.customers USING btree (account_status);


--
-- Name: idx_customers_assigned_to; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_customers_assigned_to ON public.customers USING btree (assigned_to);


--
-- Name: idx_customers_company_name; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_customers_company_name ON public.customers USING btree (company_name);


--
-- Name: idx_customers_created_by; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_customers_created_by ON public.customers USING btree (created_by);


--
-- Name: idx_customers_customer_code; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_customers_customer_code ON public.customers USING btree (customer_code);


--
-- Name: idx_gateways_enabled; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_gateways_enabled ON public.gateways USING btree (enabled);


--
-- Name: idx_gateways_last_seen; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_gateways_last_seen ON public.gateways USING btree (last_seen);


--
-- Name: idx_gateways_status; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_gateways_status ON public.gateways USING btree (status);


--
-- Name: idx_health_checked; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_health_checked ON public.gateway_health_logs USING btree (checked_at);


--
-- Name: idx_health_gateway; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_health_gateway ON public.gateway_health_logs USING btree (gateway_id);


--
-- Name: idx_health_status; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_health_status ON public.gateway_health_logs USING btree (status);


--
-- Name: idx_modems_imei; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_modems_imei ON public.modems USING btree (imei);


--
-- Name: idx_modems_imsi; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_modems_imsi ON public.modems USING btree (imsi);


--
-- Name: idx_modems_status; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_modems_status ON public.modems USING btree (status);


--
-- Name: idx_notification_templates_template_name; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_notification_templates_template_name ON public.notification_templates USING btree (template_name);


--
-- Name: idx_operator_rules_name; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_operator_rules_name ON public.operator_routing_rules USING btree (operator_name);


--
-- Name: idx_operator_rules_prefix; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_operator_rules_prefix ON public.operator_routing_rules USING btree (prefix_pattern);


--
-- Name: idx_payments_customer_id; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_payments_customer_id ON public.payments USING btree (customer_id);


--
-- Name: idx_payments_payment_reference; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_payments_payment_reference ON public.payments USING btree (payment_reference);


--
-- Name: idx_payments_payment_type; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_payments_payment_type ON public.payments USING btree (payment_type);


--
-- Name: idx_payments_processed_at; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_payments_processed_at ON public.payments USING btree (processed_at);


--
-- Name: idx_payments_status; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_payments_status ON public.payments USING btree (status);


--
-- Name: idx_rate_plans_effective_from; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_rate_plans_effective_from ON public.rate_plans USING btree (effective_from);


--
-- Name: idx_rate_plans_is_active; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_rate_plans_is_active ON public.rate_plans USING btree (is_active);


--
-- Name: idx_rate_plans_plan_code; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_rate_plans_plan_code ON public.rate_plans USING btree (plan_code);


--
-- Name: idx_revenue_customer; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_revenue_customer ON public.revenue_tracking USING btree (customer_id);


--
-- Name: idx_revenue_date; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_revenue_date ON public.revenue_tracking USING btree (billing_date);


--
-- Name: idx_revenue_source; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_revenue_source ON public.revenue_tracking USING btree (revenue_source);


--
-- Name: idx_routing_rules_is_active; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_routing_rules_is_active ON public.routing_rules USING btree (is_active);


--
-- Name: idx_routing_rules_prefix_pattern; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_routing_rules_prefix_pattern ON public.routing_rules USING btree (prefix_pattern);


--
-- Name: idx_routing_rules_rule_order; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_routing_rules_rule_order ON public.routing_rules USING btree (rule_order);


--
-- Name: idx_sim_cards_iccid; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_sim_cards_iccid ON public.sim_cards USING btree (iccid);


--
-- Name: idx_sim_cards_modem_id; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_sim_cards_modem_id ON public.sim_cards USING btree (modem_id);


--
-- Name: idx_sim_cards_msisdn; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_sim_cards_msisdn ON public.sim_cards USING btree (msisdn);


--
-- Name: idx_sim_cards_status; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_sim_cards_status ON public.sim_cards USING btree (status);


--
-- Name: idx_sim_pool_assignments_sim_card_id; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_sim_pool_assignments_sim_card_id ON public.sim_pool_assignments USING btree (sim_card_id);


--
-- Name: idx_sim_pool_assignments_sim_pool_id; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_sim_pool_assignments_sim_pool_id ON public.sim_pool_assignments USING btree (sim_pool_id);


--
-- Name: idx_sim_pools_pool_name; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_sim_pools_pool_name ON public.sim_pools USING btree (pool_name);


--
-- Name: idx_sip_calls_caller; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sip_calls_caller ON public.sip_calls USING btree (caller_number);


--
-- Name: idx_sip_calls_destination; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sip_calls_destination ON public.sip_calls USING btree (destination_number);


--
-- Name: idx_sip_calls_gateway; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sip_calls_gateway ON public.sip_calls USING btree (gateway_id);


--
-- Name: idx_sip_calls_time; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_sip_calls_time ON public.sip_calls USING btree (created_at);


--
-- Name: idx_system_config_category; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_system_config_category ON public.system_config USING btree (category);


--
-- Name: idx_system_config_config_key; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_system_config_config_key ON public.system_config USING btree (config_key);


--
-- Name: idx_user_notifications_created_at; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_user_notifications_created_at ON public.user_notifications USING btree (created_at);


--
-- Name: idx_user_notifications_is_read; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_user_notifications_is_read ON public.user_notifications USING btree (is_read);


--
-- Name: idx_user_notifications_user_id; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_user_notifications_user_id ON public.user_notifications USING btree (user_id);


--
-- Name: idx_user_sessions_expires_at; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_user_sessions_expires_at ON public.user_sessions USING btree (expires_at);


--
-- Name: idx_user_sessions_session_token; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_user_sessions_session_token ON public.user_sessions USING btree (session_token);


--
-- Name: idx_user_sessions_user_id; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_user_sessions_user_id ON public.user_sessions USING btree (user_id);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: idx_users_is_active; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_users_is_active ON public.users USING btree (is_active);


--
-- Name: idx_users_role; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_users_role ON public.users USING btree (role);


--
-- Name: idx_users_username; Type: INDEX; Schema: public; Owner: e173_user
--

CREATE INDEX idx_users_username ON public.users USING btree (username);


--
-- Name: idx_whatsapp_checked; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_whatsapp_checked ON public.whatsapp_validation_cache USING btree (checked_at);


--
-- Name: idx_whatsapp_expires; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_whatsapp_expires ON public.whatsapp_validation_cache USING btree (expires_at);


--
-- Name: blacklist set_blacklist_updated_at; Type: TRIGGER; Schema: public; Owner: e173_user
--

CREATE TRIGGER set_blacklist_updated_at BEFORE UPDATE ON public.blacklist FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: call_detail_records set_cdr_updated_at; Type: TRIGGER; Schema: public; Owner: e173_user
--

CREATE TRIGGER set_cdr_updated_at BEFORE UPDATE ON public.call_detail_records FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: customers set_customers_updated_at; Type: TRIGGER; Schema: public; Owner: e173_user
--

CREATE TRIGGER set_customers_updated_at BEFORE UPDATE ON public.customers FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: modems set_modems_updated_at; Type: TRIGGER; Schema: public; Owner: e173_user
--

CREATE TRIGGER set_modems_updated_at BEFORE UPDATE ON public.modems FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: notification_templates set_notification_templates_updated_at; Type: TRIGGER; Schema: public; Owner: e173_user
--

CREATE TRIGGER set_notification_templates_updated_at BEFORE UPDATE ON public.notification_templates FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: payments set_payments_updated_at; Type: TRIGGER; Schema: public; Owner: e173_user
--

CREATE TRIGGER set_payments_updated_at BEFORE UPDATE ON public.payments FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: rate_plans set_rate_plans_updated_at; Type: TRIGGER; Schema: public; Owner: e173_user
--

CREATE TRIGGER set_rate_plans_updated_at BEFORE UPDATE ON public.rate_plans FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: routing_rules set_routing_rules_updated_at; Type: TRIGGER; Schema: public; Owner: e173_user
--

CREATE TRIGGER set_routing_rules_updated_at BEFORE UPDATE ON public.routing_rules FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: sim_cards set_sim_cards_updated_at; Type: TRIGGER; Schema: public; Owner: e173_user
--

CREATE TRIGGER set_sim_cards_updated_at BEFORE UPDATE ON public.sim_cards FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: sim_pools set_sim_pools_updated_at; Type: TRIGGER; Schema: public; Owner: e173_user
--

CREATE TRIGGER set_sim_pools_updated_at BEFORE UPDATE ON public.sim_pools FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: system_config set_system_config_updated_at; Type: TRIGGER; Schema: public; Owner: e173_user
--

CREATE TRIGGER set_system_config_updated_at BEFORE UPDATE ON public.system_config FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: users set_users_updated_at; Type: TRIGGER; Schema: public; Owner: e173_user
--

CREATE TRIGGER set_users_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: gateways update_gateways_updated_at; Type: TRIGGER; Schema: public; Owner: e173_user
--

CREATE TRIGGER update_gateways_updated_at BEFORE UPDATE ON public.gateways FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: ai_voice_sessions ai_voice_sessions_agent_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_voice_sessions
    ADD CONSTRAINT ai_voice_sessions_agent_id_fkey FOREIGN KEY (agent_id) REFERENCES public.ai_voice_agents(id);


--
-- Name: ai_voice_sessions ai_voice_sessions_sip_call_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ai_voice_sessions
    ADD CONSTRAINT ai_voice_sessions_sip_call_id_fkey FOREIGN KEY (sip_call_id) REFERENCES public.sip_calls(id);


--
-- Name: audit_logs audit_logs_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.audit_logs
    ADD CONSTRAINT audit_logs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: blacklist blacklist_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.blacklist
    ADD CONSTRAINT blacklist_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: call_detail_records call_detail_records_modem_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.call_detail_records
    ADD CONSTRAINT call_detail_records_modem_id_fkey FOREIGN KEY (modem_id) REFERENCES public.modems(id) ON DELETE SET NULL;


--
-- Name: call_detail_records call_detail_records_sim_card_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.call_detail_records
    ADD CONSTRAINT call_detail_records_sim_card_id_fkey FOREIGN KEY (sim_card_id) REFERENCES public.sim_cards(id) ON DELETE SET NULL;


--
-- Name: customer_rate_plans customer_rate_plans_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.customer_rate_plans
    ADD CONSTRAINT customer_rate_plans_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: customer_rate_plans customer_rate_plans_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.customer_rate_plans
    ADD CONSTRAINT customer_rate_plans_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.customers(id) ON DELETE CASCADE;


--
-- Name: customer_rate_plans customer_rate_plans_rate_plan_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.customer_rate_plans
    ADD CONSTRAINT customer_rate_plans_rate_plan_id_fkey FOREIGN KEY (rate_plan_id) REFERENCES public.rate_plans(id) ON DELETE CASCADE;


--
-- Name: customers customers_assigned_to_fkey; Type: FK CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_assigned_to_fkey FOREIGN KEY (assigned_to) REFERENCES public.users(id);


--
-- Name: customers customers_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.customers
    ADD CONSTRAINT customers_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: gateway_health_logs gateway_health_logs_gateway_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.gateway_health_logs
    ADD CONSTRAINT gateway_health_logs_gateway_id_fkey FOREIGN KEY (gateway_id) REFERENCES public.gateways(id);


--
-- Name: notification_templates notification_templates_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.notification_templates
    ADD CONSTRAINT notification_templates_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: operator_routing_rules operator_routing_rules_preferred_gateway_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.operator_routing_rules
    ADD CONSTRAINT operator_routing_rules_preferred_gateway_id_fkey FOREIGN KEY (preferred_gateway_id) REFERENCES public.gateways(id);


--
-- Name: payments payments_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.customers(id) ON DELETE CASCADE;


--
-- Name: payments payments_processed_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_processed_by_fkey FOREIGN KEY (processed_by) REFERENCES public.users(id);


--
-- Name: rate_plans rate_plans_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.rate_plans
    ADD CONSTRAINT rate_plans_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: revenue_tracking revenue_tracking_ai_agent_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.revenue_tracking
    ADD CONSTRAINT revenue_tracking_ai_agent_id_fkey FOREIGN KEY (ai_agent_id) REFERENCES public.ai_voice_agents(id);


--
-- Name: revenue_tracking revenue_tracking_customer_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.revenue_tracking
    ADD CONSTRAINT revenue_tracking_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.customers(id);


--
-- Name: revenue_tracking revenue_tracking_gateway_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.revenue_tracking
    ADD CONSTRAINT revenue_tracking_gateway_id_fkey FOREIGN KEY (gateway_id) REFERENCES public.gateways(id);


--
-- Name: routing_rules routing_rules_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.routing_rules
    ADD CONSTRAINT routing_rules_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: routing_rules routing_rules_route_to_modem_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.routing_rules
    ADD CONSTRAINT routing_rules_route_to_modem_id_fkey FOREIGN KEY (route_to_modem_id) REFERENCES public.modems(id) ON DELETE SET NULL;


--
-- Name: sim_cards sim_cards_modem_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.sim_cards
    ADD CONSTRAINT sim_cards_modem_id_fkey FOREIGN KEY (modem_id) REFERENCES public.modems(id) ON DELETE SET NULL;


--
-- Name: sim_pool_assignments sim_pool_assignments_assigned_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.sim_pool_assignments
    ADD CONSTRAINT sim_pool_assignments_assigned_by_fkey FOREIGN KEY (assigned_by) REFERENCES public.users(id);


--
-- Name: sim_pool_assignments sim_pool_assignments_sim_card_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.sim_pool_assignments
    ADD CONSTRAINT sim_pool_assignments_sim_card_id_fkey FOREIGN KEY (sim_card_id) REFERENCES public.sim_cards(id) ON DELETE CASCADE;


--
-- Name: sim_pool_assignments sim_pool_assignments_sim_pool_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.sim_pool_assignments
    ADD CONSTRAINT sim_pool_assignments_sim_pool_id_fkey FOREIGN KEY (sim_pool_id) REFERENCES public.sim_pools(id) ON DELETE CASCADE;


--
-- Name: sim_pools sim_pools_created_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.sim_pools
    ADD CONSTRAINT sim_pools_created_by_fkey FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: sip_calls sip_calls_gateway_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.sip_calls
    ADD CONSTRAINT sip_calls_gateway_id_fkey FOREIGN KEY (gateway_id) REFERENCES public.gateways(id);


--
-- Name: system_config system_config_updated_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.system_config
    ADD CONSTRAINT system_config_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES public.users(id);


--
-- Name: user_notifications user_notifications_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.user_notifications
    ADD CONSTRAINT user_notifications_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: user_sessions user_sessions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: e173_user
--

ALTER TABLE ONLY public.user_sessions
    ADD CONSTRAINT user_sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--


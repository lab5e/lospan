CREATE TABLE IF NOT EXISTS lora_applications (
    eui         BIGINT       NOT NULL,
    tag         VARCHAR(128) NOT NULL,
    CONSTRAINT lora_application_pk PRIMARY KEY (eui)
);


CREATE TABLE IF NOT EXISTS lora_devices (
    eui             BIGINT       NOT NULL,
    dev_addr        CHAR(8)      NOT NULL,
    app_key         CHAR(32)     NOT NULL,
    apps_key        CHAR(32)     NOT NULL,
    nwks_key        CHAR(32)     NOT NULL,
    application_eui BIGINT       NOT NULL REFERENCES lora_application(eui),
    state           SMALLINT     NOT NULL,
    fcnt_up         INTEGER      NOT NULL DEFAULT 0,
    fcnt_dn         INTEGER      NOT NULL DEFAULT 0,
    relaxed_counter BOOLEAN      NOT NULL DEFAULT false,
    key_warning     BOOLEAN      NOT NULL DEFAULT false,
    tag             VARCHAR(128) NOT NULL,
    CONSTRAINT lora_device_pk PRIMARY KEY (eui)
);

CREATE INDEX IF NOT EXISTS lora_device_application_eui ON lora_devices(application_eui);
CREATE INDEX IF NOT EXISTS lora_device_dev_addr ON lora_devices(dev_addr);
CREATE INDEX IF NOT EXISTS lora_device_state ON lora_devices(state);


CREATE TABLE IF NOT EXISTS lora_device_nonces (
    device_eui BIGINT NOT NULL REFERENCES lora_device (eui) ON DELETE CASCADE,
    nonce      INT    NOT NULL,

    CONSTRAINT lora_device_nonce_pk PRIMARY KEY(device_eui, nonce)
);


CREATE TABLE IF NOT EXISTS lora_upstream_messages (
    device_eui      BIGINT        NOT NULL REFERENCES lora_device (eui) ON DELETE CASCADE, 
    data            VARCHAR(512)  NOT NULL, 
    time_stamp      BIGINT        NOT NULL, 
    gateway_eui     BIGINT        NOT NULL, 
    rssi            INTEGER       NOT NULL,
    snr             NUMERIC(6,3)  NOT NULL,
    frequency       NUMERIC(6,3)  NOT NULL,
    data_rate       VARCHAR(20)   NOT NULL,
    dev_addr        CHAR(8)       NOT NULL,

    CONSTRAINT lora_device_data_pk PRIMARY KEY(device_eui, time_stamp)
);

CREATE INDEX IF NOT EXISTS lora_device_data_device_eui ON lora_upstream_messages(device_eui);


CREATE TABLE IF NOT EXISTS lora_sequences (
    identifier VARCHAR(128) NOT NULL, 
    counter    BIGINT       NOT NULL, 

    CONSTRAINT lora_sequence_pk PRIMARY KEY (identifier)
);

CREATE INDEX IF NOT EXISTS lora_sequence_identifier ON lora_sequences(identifier);


CREATE TABLE IF NOT EXISTS lora_gateways (
    gateway_eui BIGINT     NOT NULL,
    latitude    NUMERIC(12,8) NULL,
    longitude   NUMERIC(12,8) NULL,
    altitude    NUMERIC(8,3)  NULL,
    ip          VARCHAR(64)   NOT NULL,
    strict_ip   BOOL          NOT NULL,

    CONSTRAINT lora_gateway_pk PRIMARY KEY (gateway_eui)
);


CREATE TABLE IF NOT EXISTS lora_downstream_messages (
    device_eui   BIGINT NOT NULL REFERENCES lora_device(eui) ON DELETE CASCADE,
    data         VARCHAR(256) NOT NULL,
    port         INTEGER NOT NULL,
    ack          BOOLEAN NOT NULL DEFAULT false,
    created_time INTEGER NOT NULL,
    sent_time    INTEGER DEFAULT 0,
    ack_time     INTEGER DEFAULT 0,

    CONSTRAINT lora_downstream_message_pk PRIMARY KEY (device_eui, created_time)
);

CREATE INDEX IF NOT EXISTS lora_downstream_messages_created ON lora_downstream_messages(created_time);  
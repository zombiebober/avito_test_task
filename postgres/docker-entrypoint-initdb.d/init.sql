CREATE TABLE advert
(
    ID SERIAL PRIMARY KEY,
    title text ,
    Description text,
    Price numeric,
    Photo_link text[],
    time_create timestamp with time zone DEFAULT current_timestamp
);
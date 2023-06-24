-- Active: 1685868021088@@localhost@3306@simple_bank
SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;
begin;
INSERT INTO transfers
    (
    from_account_id,
    to_account_id,
    amount
    )
VALUES
    (13, 14, 10);
-- SELECT SLEEP(30);
-- update accounts set balance = balance - 10 where id = 13;
-- update accounts set balance = balance + 10 where id = 14;
SELECT TRX_MYSQL_THREAD_ID,dl.*, it.trx_query
FROM performance_schema.data_locks dl
JOIN information_schema.innodb_trx it ON dl.engine_transaction_id = it.trx_id
where TRX_MYSQL_THREAD_ID = CONNECTION_ID();
-- commit;
ROLLBACK;
CREATE TABLE [accounts] (
  [id] BIGINT PRIMARY KEY,
  [owner] nvarchar(255) NOT NULL,
  [balance] decimal NOT NULL,
  [currency] nvarchar(255) NOT NULL,
  [created_at] datetimeoffset NOT NULL DEFAULT (getdate())
)
GO

CREATE TABLE [entries] (
  [id] BIGINT PRIMARY KEY,
  [account_id] BIGINT,
  [amount] BIGINT NOT NULL,
  [created_at] datetimeoffset NOT NULL DEFAULT (getdate())
)
GO

CREATE TABLE [transfers] (
  [id] BIGINT PRIMARY KEY,
  [from_account_id] BIGINT NOT NULL,
  [to_account_id] BIGINT NOT NULL,
  [amount] BIGINT NOT NULL,
  [created_at] datetimeoffset NOT NULL DEFAULT (getdate())
)
GO

CREATE INDEX [accounts_index_0] ON [accounts] ("owner")
GO

CREATE INDEX [entries_index_1] ON [entries] ("account_id")
GO

CREATE INDEX [transfers_index_2] ON [transfers] ("from_account_id")
GO

CREATE INDEX [transfers_index_3] ON [transfers] ("to_account_id")
GO

CREATE INDEX [transfers_index_4] ON [transfers] ("from_account_id", "to_account_id")
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '可以是正或者负，表示存入或者取出',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'entries',
@level2type = N'Column', @level2name = 'amount';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Content of the post',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'transfers',
@level2type = N'Column', @level2name = 'to_account_id';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '只能是正数',
@level0type = N'Schema', @level0name = 'dbo',
@level1type = N'Table',  @level1name = 'transfers',
@level2type = N'Column', @level2name = 'amount';
GO

ALTER TABLE [entries] ADD FOREIGN KEY ([account_id]) REFERENCES [accounts] ([id])
GO

ALTER TABLE [transfers] ADD FOREIGN KEY ([from_account_id]) REFERENCES [accounts] ([id])
GO

ALTER TABLE [transfers] ADD FOREIGN KEY ([to_account_id]) REFERENCES [accounts] ([id])
GO

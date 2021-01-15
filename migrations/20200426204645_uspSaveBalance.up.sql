create or alter proc dbo.uspSaveBalance @AccountId int,
                                        @AvailableBalance decimal(9, 5),
                                        @OutstandingTransactions decimal(9, 5),
                                        @GivenCredit decimal(9, 5),
                                        @Currency varchar(50) as
begin
    set nocount on

    MERGE dbo.Balance AS t
    USING (select @AccountId, @AvailableBalance, @OutstandingTransactions, @GivenCredit, @Currency) s
        (AccountId, AvailableBalance, OutstandingTransactions, GivenCredit, Currency)
    on s.AccountId = t.AccountId

    WHEN MATCHED THEN
        UPDATE
        SET t.AvailableBalance        = s.AvailableBalance,
            t.OutstandingTransactions = s.OutstandingTransactions,
            t.GivenCredit             = s.GivenCredit,
            t.Currency                = s.Currency,
            t.UpdatedAt               = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (AccountId,
                AvailableBalance,
                OutstandingTransactions,
                GivenCredit,
                Currency)
        VALUES (s.AccountId,
                s.AvailableBalance,
                s.OutstandingTransactions,
                s.GivenCredit,
                s.Currency);
end

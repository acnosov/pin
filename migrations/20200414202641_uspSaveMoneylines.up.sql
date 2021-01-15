create or alter proc dbo.uspSaveMoneylines @TVP dbo.MoneylineType READONLY as
begin
    set nocount on

    MERGE dbo.Moneyline AS t
    USING @TVP s
    ON (t.LineId = s.LineId)

    WHEN MATCHED THEN
        UPDATE
        SET Home      = s.Home,
            Away      = s.Away,
            Draw      = s.Draw,
            UpdatedAt = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (LineId,
                Home,
                Away,
                Draw)
        VALUES (s.LineId,
                s.Home,
                s.Away,
                s.Draw);
end

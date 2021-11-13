create or alter proc dbo.uspSavePeriods @SportId int, @TVP dbo.PeriodType READONLY as
begin
    set nocount on

    MERGE dbo.Period AS t
    USING @TVP s
    ON (t.Number = s.Number and t.SportId = @SportId)

    WHEN MATCHED THEN
        UPDATE
        SET Description                = s.Description,
            ShortDescription           = s.ShortDescription,
            SpreadDescription          = s.SpreadDescription,
            MoneylineDescription       = s.MoneylineDescription,
            TotalDescription           = s.TotalDescription,
            Team1TotalDescription      = s.Team1TotalDescription,
            Team2TotalDescription      = s.Team2TotalDescription,
            SpreadShortDescription     = s.SpreadShortDescription,
            MoneylineShortDescription  = s.MoneylineShortDescription,
            TotalShortDescription      = s.TotalShortDescription,
            Team1TotalShortDescription = s.Team1TotalShortDescription,
            Team2TotalShortDescription = s.Team2TotalShortDescription,
            UpdatedAt                  = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (Number, SportId, Description, ShortDescription, SpreadDescription, MoneylineDescription,
                TotalDescription, Team1TotalDescription, Team2TotalDescription, SpreadShortDescription,
                MoneylineShortDescription, TotalShortDescription, Team1TotalShortDescription,
                Team2TotalShortDescription)
        VALUES (s.Number, @SportId, s.Description, s.ShortDescription, s.SpreadDescription, s.MoneylineDescription,
                s.TotalDescription, s.Team1TotalDescription, s.Team2TotalDescription, s.SpreadShortDescription,
                s.MoneylineShortDescription, s.TotalShortDescription, s.Team1TotalShortDescription,
                s.Team2TotalShortDescription);
end
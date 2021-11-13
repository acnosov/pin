create or alter proc dbo.uspSaveSpecialEvents @TVP dbo.SpecialEventType READONLY as
begin
    set nocount on

    MERGE dbo.SpecialEvent AS t
    USING @TVP s
    ON (t.Id = s.Id and t.SpecialId = s.SpecialId)

    WHEN MATCHED THEN
        UPDATE
        SET PeriodNumber = s.PeriodNumber,
            Home         = s.Home,
            Away         = s.Away,
            UpdatedAt    = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (Id,
                PeriodNumber,
                Home,
                Away,
                SpecialId)
        VALUES (s.Id,
                s.PeriodNumber,
                s.Home,
                s.Away,
                s.SpecialId);
end

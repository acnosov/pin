create or alter proc dbo.uspSaveSpecialContestants @TVP dbo.SpecialContestantType READONLY as
begin
    set nocount on

    MERGE dbo.SpecialContestant AS t
    USING @TVP s
    ON (t.Id = s.Id)

    WHEN MATCHED THEN
        UPDATE
        SET Name      = s.Name,
            RotNum    = s.RotNum,
            SpecialId = s.SpecialId,
            UpdatedAt = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (Id,
                Name,
                RotNum,
                SpecialId)
        VALUES (s.Id,
                s.Name,
                s.RotNum,
                s.SpecialId);
end

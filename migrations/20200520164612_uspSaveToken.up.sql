create or alter proc dbo.uspSaveToken @Session varchar(100),
                                      @ApiKey varchar(100),
                                      @Device varchar(100),
                                      @TrustCode varchar(100)
as
begin
    set nocount on;
    MERGE dbo.Auth AS t
    USING (select @Session) s (Session)
    on t.Session = s.Session
    WHEN MATCHED THEN
        UPDATE set ApiKey = @ApiKey, Device = @Device, TrustCode = @TrustCode, LastCheckAt = sysdatetimeoffset()
    WHEN NOT MATCHED THEN
        INSERT (Session, ApiKey, Device, TrustCode) VALUES (@Session, @ApiKey, @Device, @TrustCode);
end
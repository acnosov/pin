create or alter proc dbo.uspCleanOldData
as
begin
    set nocount on;
    SELECT 1;
    SELECT 1;
end;
go


select count(*)
from Line
exec dbo.uspCleanOldData

truncate table dbo.Spread
truncate table dbo.Total
truncate table dbo.Moneyline
truncate table dbo.Line
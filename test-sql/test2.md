### GetStudentByID2

```sql
select * from student where id = {{. | bind}}
```

### GetStudentByID3

```sql
select * from student where id in ({{. | bind}})
```

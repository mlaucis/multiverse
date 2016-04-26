WITH
table_stats AS (
  SELECT
    psut.schemaname,
    psut.relname,
    psut.n_live_tup,
    1.0 * psut.idx_scan / greatest(1, psut.seq_scan + psut.idx_scan) AS index_use_ratio
  FROM
    pg_stat_user_tables psut
  WHERE
    psut.n_live_tup > 10000
  ORDER BY
    psut.n_live_tup DESC
),
table_io AS (
  SELECT
    psiut.relname,
    sum(psiut.heap_blks_read) AS table_page_read,
    sum(psiut.heap_blks_hit)  AS table_page_hit,
    sum(psiut.heap_blks_hit) / greatest(1, sum(psiut.heap_blks_hit) + sum(psiut.heap_blks_read)) AS table_hit_ratio
  FROM
    pg_statio_user_tables psiut
  GROUP BY
    psiut.relname
  ORDER BY
    table_page_read DESC
),
index_io AS (
  SELECT
    psiui.relname,
    psiui.indexrelname,
    sum(psiui.idx_blks_read) AS idx_page_read,
    sum(psiui.idx_blks_hit) AS idx_page_hit,
    1.0 * sum(psiui.idx_blks_hit) / greatest(1.0, sum(psiui.idx_blks_hit) + sum(psiui.idx_blks_read)) AS idx_hit_ratio
  FROM
    pg_statio_user_indexes psiui
  GROUP BY
    psiui.relname, psiui.indexrelname
  ORDER BY
    sum(psiui.idx_blks_read)
  DESC
)
SELECT
	ts.schemaname,
  ts.relname,
  ts.n_live_tup,
  ts.index_use_ratio,
  ti.table_page_read,
  ti.table_page_hit,
  ti.table_hit_ratio,
  ii.indexrelname,
  ii.idx_page_read,
  ii.idx_page_hit,
  ii.idx_hit_ratio
FROM
  table_stats ts
LEFT OUTER JOIN
  table_io ti ON ti.relname = ts.relname
LEFT OUTER JOIN
  index_io ii ON ii.relname = ts.relname
ORDER BY
  ti.table_page_read DESC,
  ii.idx_page_read DESC;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS job_queue;

CREATE TABLE projects (
                          id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                          user_id UUID NOT NULL,
                          name VARCHAR(255) NOT NULL,
                          description TEXT NOT NULL,
                          base_keywords TEXT[] NOT NULL DEFAULT '{}',
                          status VARCHAR(20) NOT NULL DEFAULT 'pending',
                          result_data JSONB DEFAULT NULL,
                          ai_result_data JSONB DEFAULT NULL,
                          seo_result_data JSONB DEFAULT NULL,
                          created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_created_at ON projects(created_at DESC);

INSERT INTO projects (user_id, name, description, base_keywords, status)
VALUES
    ('550e8400-e29b-41d4-a716-446655440000', 'Проект 1', 'Описание проекта 1', ARRAY['seo', 'оптимизация'], 'completed'),
    ('550e8400-e29b-41d4-a716-446655440001', 'Проект 2', 'Описание проекта 2', ARRAY['маркетинг', 'анализ'], 'pending')
ON CONFLICT DO NOTHING;


CREATE TABLE job_queue (
                           id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                           entity_id UUID NOT NULL,
                           stage VARCHAR(50) NOT NULL,
                           created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_job_queue_stage ON job_queue(stage);
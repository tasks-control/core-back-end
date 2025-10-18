-- +goose Up
-- +goose StatementBegin

-- =====================================================
-- Table: members (Users)
-- =====================================================
CREATE TABLE members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(50) NOT NULL UNIQUE,
    full_name VARCHAR(100),
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_members_email ON members(email);
CREATE INDEX idx_members_username ON members(username);

-- =====================================================
-- Table: boards (Boards)
-- =====================================================
CREATE TABLE boards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    name_board_unique VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    password_hash VARCHAR(255) NOT NULL,
    id_member_creator UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_boards_creator 
        FOREIGN KEY (id_member_creator) 
        REFERENCES members(id) 
        ON DELETE CASCADE,
    
    CONSTRAINT chk_name_board_unique_format 
        CHECK (name_board_unique ~ '^[a-z0-9-]+$')
);

CREATE UNIQUE INDEX idx_boards_name_unique ON boards(name_board_unique);
CREATE INDEX idx_boards_creator ON boards(id_member_creator);
CREATE INDEX idx_boards_updated_at ON boards(updated_at DESC);

-- =====================================================
-- Table: board_members (Board Membership)
-- =====================================================
CREATE TABLE board_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_board UUID NOT NULL,
    id_member UUID NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'member',
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_board_members_board 
        FOREIGN KEY (id_board) 
        REFERENCES boards(id) 
        ON DELETE CASCADE,
    
    CONSTRAINT fk_board_members_member 
        FOREIGN KEY (id_member) 
        REFERENCES members(id) 
        ON DELETE CASCADE,
    
    CONSTRAINT chk_board_members_role 
        CHECK (role IN ('owner', 'member')),
    
    CONSTRAINT uq_board_member 
        UNIQUE (id_board, id_member)
);

CREATE INDEX idx_board_members_board ON board_members(id_board);
CREATE INDEX idx_board_members_member ON board_members(id_member);
CREATE INDEX idx_board_members_role ON board_members(id_board, role);

-- =====================================================
-- Table: starred_boards (Favorite Boards)
-- =====================================================
CREATE TABLE starred_boards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_board UUID NOT NULL,
    id_member UUID NOT NULL,
    starred_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_starred_boards_board 
        FOREIGN KEY (id_board) 
        REFERENCES boards(id) 
        ON DELETE CASCADE,
    
    CONSTRAINT fk_starred_boards_member 
        FOREIGN KEY (id_member) 
        REFERENCES members(id) 
        ON DELETE CASCADE,
    
    CONSTRAINT uq_starred_board_member 
        UNIQUE (id_board, id_member)
);

CREATE INDEX idx_starred_boards_member ON starred_boards(id_member);
CREATE INDEX idx_starred_boards_board ON starred_boards(id_board);

-- =====================================================
-- Table: lists (Board Lists/Columns)
-- =====================================================
CREATE TABLE lists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    id_board UUID NOT NULL,
    position DOUBLE PRECISION NOT NULL,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_lists_board 
        FOREIGN KEY (id_board) 
        REFERENCES boards(id) 
        ON DELETE CASCADE
);

CREATE INDEX idx_lists_board ON lists(id_board);
CREATE INDEX idx_lists_board_position ON lists(id_board, position) WHERE archived = FALSE;
CREATE INDEX idx_lists_archived ON lists(id_board, archived);

-- =====================================================
-- Table: cards (Task Cards)
-- =====================================================
CREATE TABLE cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(200) NOT NULL,
    description TEXT,
    id_list UUID NOT NULL,
    position DOUBLE PRECISION NOT NULL,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_cards_list 
        FOREIGN KEY (id_list) 
        REFERENCES lists(id) 
        ON DELETE CASCADE,
    
    CONSTRAINT fk_cards_creator 
        FOREIGN KEY (created_by) 
        REFERENCES members(id) 
        ON DELETE CASCADE
);

CREATE INDEX idx_cards_list ON cards(id_list);
CREATE INDEX idx_cards_list_position ON cards(id_list, position) WHERE archived = FALSE;
CREATE INDEX idx_cards_archived ON cards(id_list, archived);
CREATE INDEX idx_cards_creator ON cards(created_by);
CREATE INDEX idx_cards_updated_at ON cards(updated_at DESC);

-- =====================================================
-- Table: refresh_tokens (JWT Refresh Tokens)
-- =====================================================
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id_member UUID NOT NULL,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    
    CONSTRAINT fk_refresh_tokens_member 
        FOREIGN KEY (id_member) 
        REFERENCES members(id) 
        ON DELETE CASCADE
);

CREATE INDEX idx_refresh_tokens_member ON refresh_tokens(id_member);
CREATE INDEX idx_refresh_tokens_hash ON refresh_tokens(token_hash) WHERE revoked = FALSE;
CREATE INDEX idx_refresh_tokens_expires ON refresh_tokens(expires_at) WHERE revoked = FALSE;

-- =====================================================
-- Functions and Triggers
-- =====================================================

-- Function to automatically update updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers for updated_at
CREATE TRIGGER update_members_updated_at 
    BEFORE UPDATE ON members
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_boards_updated_at 
    BEFORE UPDATE ON boards
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_lists_updated_at 
    BEFORE UPDATE ON lists
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_cards_updated_at 
    BEFORE UPDATE ON cards
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop triggers
DROP TRIGGER IF EXISTS update_cards_updated_at ON cards;
DROP TRIGGER IF EXISTS update_lists_updated_at ON lists;
DROP TRIGGER IF EXISTS update_boards_updated_at ON boards;
DROP TRIGGER IF EXISTS update_members_updated_at ON members;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order (respecting foreign keys)
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS cards;
DROP TABLE IF EXISTS lists;
DROP TABLE IF EXISTS starred_boards;
DROP TABLE IF EXISTS board_members;
DROP TABLE IF EXISTS boards;
DROP TABLE IF EXISTS members;

-- +goose StatementEnd

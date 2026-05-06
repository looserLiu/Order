# Order - Accounting & Inventory Management App

## Project Overview

**Project Name:** Order
**Type:** Flutter Mobile Application (Local-First)
**Core Functionality:** Dual-module application combining personal/small-business accounting (transactions, accounts, budgets, smart categorization) with inventory management (products, warehouses, stock flows, expiration tracking)
**Target Users:** Small business owners, freelancers, and individuals needing integrated accounting and inventory tracking

## Architecture

### Technology Stack

- **Framework:** Flutter 3.x (Dart)
- **Database:** SQLite (sqflite package) for local-first offline capability
- **State Management:** Riverpod 2.x (providers)
- **Architecture Pattern:** Clean Architecture (Presentation / Domain / Data layers)

### Project Structure

```
lib/
  core/
    constants/        # App constants (colors, strings, dimensions)
    utils/            # Utility functions (formatters, validators)
    theme/            # App theme configuration
    widgets/          # Reusable widgets (buttons, cards, inputs)
  data/
    database/         # SQLite database helper, migrations
    models/           # Data models (database row mapping)
    repositories/     # Repository implementations
  domain/
    entities/         # Business entities
    repositories/     # Repository interfaces
    usecases/         # Business logic use cases
  presentation/
    pages/            # Screen pages
    widgets/           # Page-specific widgets
    providers/        # Riverpod providers
assets/
  icons/               # App icons
test/
  unit/               # Unit tests
  widget/             # Widget tests
```

## Data Models

### Accounting Module

#### Account
| Field | Type | Description |
|-------|------|-------------|
| id | INTEGER (PK) | Auto-increment primary key |
| name | TEXT | Account name |
| type | TEXT | 'cash', 'bank', 'credit_card', 'digital' |
| balance | REAL | Current balance |
| currency | TEXT | Currency code (default: CNY) |
| icon | TEXT | Icon identifier |
| color | INTEGER | Color value |
| created_at | INTEGER | Unix timestamp |
| updated_at | INTEGER | Unix timestamp |

#### Transaction
| Field | Type | Description |
|-------|------|-------------|
| id | INTEGER (PK) | Auto-increment primary key |
| account_id | INTEGER (FK) | Reference to Account |
| category_id | INTEGER (FK) | Reference to Category |
| amount | REAL | Transaction amount |
| type | TEXT | 'income' or 'expense' |
| description | TEXT | Transaction description |
| date | INTEGER | Transaction date (unix timestamp) |
| created_at | INTEGER | Unix timestamp |

#### Category
| Field | Type | Description |
|-------|------|-------------|
| id | INTEGER (PK) | Auto-increment primary key |
| name | TEXT | Category name |
| type | TEXT | 'income' or 'expense' |
| icon | TEXT | Icon identifier |
| color | INTEGER | Color value |
| parent_id | INTEGER (FK, nullable) | Parent category for subcategories |
| is_system | INTEGER | 1 if system default, 0 if user-created |
| usage_count | INTEGER | Times category was used (for smart suggestions) |

#### Budget
| Field | Type | Description |
|-------|------|-------------|
| id | INTEGER (PK) | Auto-increment primary key |
| category_id | INTEGER (FK) | Reference to Category |
| amount | REAL | Budget limit |
| period | TEXT | 'monthly', 'weekly', 'yearly' |
| start_date | INTEGER | Budget period start |
| end_date | INTEGER | Budget period end |

### Inventory Module

#### Product
| Field | Type | Description |
|-------|------|-------------|
| id | INTEGER (PK) | Auto-increment primary key |
| name | TEXT | Product name |
| sku | TEXT | Stock keeping unit code |
| category | TEXT | Product category |
| unit | TEXT | Unit of measure (pcs, kg, box) |
| cost_price | REAL | Cost price |
| sale_price | REAL | Sale price |
| image_url | TEXT | Product image path |
| created_at | INTEGER | Unix timestamp |
| updated_at | INTEGER | Unix timestamp |

#### Warehouse
| Field | Type | Description |
|-------|------|-------------|
| id | INTEGER (PK) | Auto-increment primary key |
| name | TEXT | Warehouse name |
| location | TEXT | Physical location |
| description | TEXT | Description |
| is_active | INTEGER | 1 if active, 0 if inactive |
| created_at | INTEGER | Unix timestamp |

#### InventoryFlow
| Field | Type | Description |
|-------|------|-------------|
| id | INTEGER (PK) | Auto-increment primary key |
| product_id | INTEGER (FK) | Reference to Product |
| warehouse_id | INTEGER (FK) | Reference to Warehouse |
| flow_type | TEXT | 'in', 'out', 'transfer', 'adjust' |
| quantity | REAL | Flow quantity |
| batch_number | TEXT | Batch number for tracking |
| expiration_date | INTEGER | Expiration date (unix timestamp, nullable) |
| cost_at_flow | REAL | Cost at time of flow |
| reference_id | TEXT | Reference to related transaction/order |
| date | INTEGER | Flow date (unix timestamp) |
| created_at | INTEGER | Unix timestamp |

#### CostCategory
| Field | Type | Description |
|-------|------|-------------|
| id | INTEGER (PK) | Auto-increment primary key |
| name | TEXT | Cost category name |
| type | TEXT | 'purchase', 'storage', 'transport', 'other' |
| description | TEXT | Description |

## Feature Modules

### Tab 1: Accounting (记账)

**Pages:**
- **Account List** - View all accounts with balances
- **Account Detail** - View transactions for specific account
- **Transaction List** - All transactions with filters (date range, account, category)
- **Add/Edit Transaction** - Create or modify transaction
- **Category Management** - Manage income/expense categories

**Features:**
- Multi-account support (cash, bank, credit card, digital wallet)
- Income/expense recording
- Transaction search and filtering
- Account balance overview
- Quick add transaction via floating action button

### Tab 2: Inventory (存库)

**Pages:**
- **Product List** - View all products with stock levels
- **Product Detail** - View product details and stock by warehouse
- **Add/Edit Product** - Create or modify product
- **Warehouse List** - View all warehouses
- **Add/Edit Warehouse** - Create or modify warehouse
- **Stock Flow List** - View all inventory movements
- **Add/Edit Stock Flow** - Record stock in/out/transfer
- **Expiration Alert** - Products expiring soon

**Features:**
- Multi-warehouse inventory tracking
- Batch and expiration date tracking
- Stock flow recording (in/out/transfer/adjust)
- Cost category tracking for inventory flows
- Low stock alerts
- Expiration date alerts

### Tab 3: Reports (报表)

**Pages:**
- **Dashboard** - Overview with charts and summaries
- **Income/Expense Report** - Monthly income vs expense analysis
- **Budget Report** - Budget vs actual spending
- **Inventory Report** - Stock valuation and turnover

**Features:**
- Visual charts (bar, pie, line)
- Date range filtering
- Export reports (future enhancement)

### Tab 4: Profile (我的)

**Pages:**
- **Profile Settings** - User settings and preferences
- **Data Management** - Backup/restore data
- **About** - App information

**Features:**
- Currency settings
- Data backup/restore
- Category management access
- Smart categorization toggle

## Smart Categorization Logic

The smart categorization system learns from transaction history to suggest categories for new transactions.

### Algorithm:

1. **Keyword Matching**
   - Extract keywords from transaction description
   - Match against historical transactions with same keywords
   - Return category with highest keyword match score

2. **Usage-Based Suggestion**
   - Track category usage count
   - When amount is similar (±10%) to a previous transaction, suggest same category
   - Rank by recency: recent transactions weighted higher

3. **Time-Based Pattern**
   - Detect recurring transactions (rent, subscriptions)
   - Suggest same category if same day of month (±3 days)

### Suggestion Display:
- Show top 3 suggested categories when adding transaction
- Highlight most likely category based on description input
- User can still select any category or create new one

## Page Navigation

Bottom navigation with 4 tabs:

```
┌─────────────────────────────────────────┐
│  [记账]    [存库]    [报表]    [我的]   │
│ Account  Inventory  Reports  Profile   │
└─────────────────────────────────────────┘
```

- **记账 (Account)** - Default landing tab
- **存库 (Inventory)** - Products and warehouse management
- **报表 (Reports)** - Charts and financial summaries
- **我的 (Profile)** - Settings and data management

## Database Schema Summary

```
┌──────────────┐     ┌──────────────┐
│   Account    │     │  Category    │
├──────────────┤     ├──────────────┤
│ Transaction  │────<│             │
├──────────────┤     │             │
│   Budget     │────<│             │
└──────────────┘     └──────────────┘

┌──────────────┐     ┌──────────────┐
│   Product    │──────│CostCategory │
├──────────────┤     └──────────────┘
│   Warehouse  │
├──────────────┤
│ InventoryFlow│────>│  references  │
└──────────────┘     │  Transaction │
                     └──────────────┘
```

## Implementation Notes

- All timestamps stored as Unix epoch integers
- Currency amounts stored as REAL (double precision)
- Soft delete not implemented in v1 (hard delete only)
- Database migrations handled via version number in database helper
- Initial seed data: default income/expense categories
export interface Group {
  id: string;
  name: string;
  description: string;
  created_at: string;
}

export interface User {
  id: string;
  name: string;
  surname: string;
  email: string;
}

export interface ExpenseSplit {
  expense_id: string;
  user_id: string;
  amount: string;
}

export interface Expense {
  id: string;
  group_id: string;
  payer_id: string;
  amount: string;
  currency: string;
  description: string;
  created_at: string;
  splits?: ExpenseSplit[];
}

export interface CreateExpenseSplitRequest {
  user_id: string;
  amount: string;
}

export interface CreateExpenseRequest {
  payer_id: string;
  amount: string;
  currency: string;
  description: string;
  splits: CreateExpenseSplitRequest[];
}

export interface MemberBalance {
  user: User;
  balance: string;
}
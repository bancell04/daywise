export interface Task {
    id?: number; 
    title: string;
    category: number;
    start?: string;
    end?: string;
}

export interface Category {
    id?: number;
    name: string;
    color: string;
}
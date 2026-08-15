export type Settings = Record<string, string>;

export interface NavItem {
  id: number;
  parent_id: number;
  label: string;
  path: string;
  location: "header" | "footer";
  open_new: boolean;
  children?: NavItem[];
}

export interface Category {
  id: number;
  parent_id: number;
  name: string;
  slug: string;
  image_url: string;
  summary: string;
  home_featured: boolean;
  children?: Category[];
}

export interface ProductImage {
  id: number;
  image_url: string;
  alt_text: string;
}

export interface Product {
  id: number;
  category_id: number;
  name: string;
  slug: string;
  sku: string;
  summary: string;
  description: string;
  main_image: string;
  video_url: string;
  material: string;
  size: string;
  thread_standard: string;
  pressure_rating: string;
  temperature_range: string;
  moq: string;
  standard: string;
  application: string;
  hot_tags: string;
  category: Category;
  images: ProductImage[];
}

export interface Article {
  id: number;
  content_type: "news" | "knowledge" | "blog";
  category: string;
  title: string;
  slug: string;
  summary: string;
  content: string;
  cover_image: string;
  published_at: string;
}

export interface PageData {
  title: string;
  slug: string;
  subtitle: string;
  content: string;
  cover_image: string;
  template: string;
}

export interface Banner { id: number; title: string; subtitle: string; image_url: string; link_url: string; }
export interface Benefit { id: number; title: string; icon: string; }
export interface Partner { id: number; name: string; logo_url: string; website_url: string; }
export interface TeamMember { id: number; name: string; role: string; email: string; image_url: string; bio: string; }
export interface Industry { id: number; title: string; subtitle: string; description: string; image_url: string; link_url: string; }
export interface Video { id: number; title: string; description: string; cover_url: string; video_url: string; }

export interface HomeData {
  settings: Settings;
  banners: Banner[];
  benefits: Benefit[];
  featured_categories: Category[];
  hot_products: Product[];
  partners: Partner[];
  team: TeamMember[];
  industries: Industry[];
  news: Article[];
}

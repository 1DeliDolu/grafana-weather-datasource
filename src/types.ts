import { DataSourceJsonData, SelectableValue } from '@grafana/data';
import { DataQuery,  } from '@grafana/schema';

export interface MyQuery extends DataQuery {
  city?: string;
}

export const DEFAULT_QUERY: Partial<MyQuery> = {
  city: "Marburg",
};

export interface DataPoint {
  Time: number;
  Value: number;
}

export interface DataSourceResponse {
  datapoints: DataPoint[];
}

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {
  path?: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface MySecureJsonData {
  apiKey?: string;
}

//städte in deutsclande Datenquellen

export const stadte: SelectableValue[] = [
  { label: "Aachen", value: "Aachen" },
  { label: "Augsburg", value: "Augsburg" },
  { label: "Berlin", value: "Berlin" },
  { label: "Bremen", value: "Bremen" },
  { label: "Dresden", value: "Dresden" },
  { label: "Dusseldorf", value: "Dusseldorf" },
  { label: "Frankfurt", value: "Frankfurt" },
  { label: "Hamburg", value: "Hamburg" },
  { label: "Hannover", value: "Hannover" },
  { label: "Köln", value: "Köln" },
  { label: "Leipzig", value: "Leipzig" },
  { label: "Mannheim", value: "Mannheim" },
  { label: "München", value: "München" },
  { label: "Nürnberg", value: "Nürnberg" },
  { label: "Stuttgart", value: "Stuttgart" },
  { label: "Würzburg", value: "Würzburg" },
  { label: "Zwickau", value: "Zwickau" },
  { label: "Marburg", value: "Marburg" },
  { label: "Wuppertal", value: "Wuppertal" },
  { label: "Dortmund", value: "Dortmund" },
  { label: "Bielefeld", value: "Bielefeld" },
];
import pandas as pd

# Load the Excel file
file_path = '1 KK_ART Pondokrejo.xlsx'
df = pd.read_excel(file_path, nrows=1)

# Print columns with their indices
for i, col in enumerate(df.columns):
    print(f"{i}: {col}")

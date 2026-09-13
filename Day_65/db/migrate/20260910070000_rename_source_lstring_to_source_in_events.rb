class RenameSourceLstringToSourceInEvents < ActiveRecord::Migration[8.1]
  def change
    rename_column :events, :sourceLstring, :source
  end
end

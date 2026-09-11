class CreateEvents < ActiveRecord::Migration[8.1]
  def change
    create_table :events do |t|
      t.string :level
      t.text :message
      t.string :source
      t.integer :user_id

      t.timestamps
    end
  end
end

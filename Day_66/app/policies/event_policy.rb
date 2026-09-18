class EventPolicy < ApplicationPolicy
  def index?
    true
  end

  def show?
    admin? || owner?
  end

  def create?
    user.present?
  end

  def update?
    admin? || owner?
  end

  def destroy?
    admin? || owner?
  end

  private

  def admin?
    user.admin?
  end

  def owner?
    record.user_id == user.id
  end

  class Scope < ApplicationPolicy::Scope
    def resolve
      if user.admin?
        scope.all
      else
        scope.where(user_id: user.id)
      end
    end
  end
end

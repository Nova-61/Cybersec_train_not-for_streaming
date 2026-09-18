class EventsController < ApplicationController
  before_action :authenticate_user!
  before_action :set_event, only: [:show, :edit, :update, :destroy]

  # GET /events
  def index
    @events = policy_scope(Event).order(created_at: :desc)
  end

  # GET /events/:id
  def show
    authorize @event
  end

  # GET /events/new
  def new
    @event = Event.new
  end

  # POST /events
  def create
    @event = Event.new(event_params)
    @event.user_id = current_user.id
    authorize @event

    if @event.save
      redirect_to @event, notice: "Событие создано!"
    else
      render :new, status: :unprocessable_entity
    end
  end

  # GET /events/:id/edit
  def edit
    authorize @event
  end

  # PATCH/PUT /events/:id
  def update
    authorize @event

    if @event.update(event_params)
      redirect_to @event, notice: "Событие обновлено!"
    else
      render :edit, status: :unprocessable_entity
    end
  end

  # DELETE /events/:id
  def destroy
    authorize @event
    @event.destroy
    redirect_to events_path, notice: "Событие удалено!", status: :see_other
  end

  private

  def set_event
    @event = Event.find(params[:id])
  end

  def event_params
    params.require(:event).permit(:level, :message, :source)
  end
end
